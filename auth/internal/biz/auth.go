package biz

import (
	"context"
	"errors"
	"net/mail"

	kerrors "github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"

	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/utils"
	"github.com/databonfire/bonfire/auth/user"
)

type AuthUsecase struct {
	conf     *conf.Biz
	userRepo UserRepo
	// hooks
	hooks map[string]user.HookFunc
}

func NewAuthUsecase(c *conf.Biz, userRepo UserRepo, hooks map[string]user.HookFunc) *AuthUsecase {
	return &AuthUsecase{
		c,
		userRepo,
		hooks,
	}
}

func (au *AuthUsecase) Register(ctx context.Context, req *pb.RegisterRequest) error {
	if len(req.Phone) == 0 && len(req.Email) == 0 {
		// Phone and email cannot both be empty
		errMsg := map[string]string{
			"email": "email is empty",
			"phone": "phone is empty",
		}
		return kerrors.BadRequest("phone and email cannot both be empty", "").WithMetadata(errMsg)
	}

	_u, err := au.userRepo.Find(ctx, req.Name, req.Email, req.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	// name, email, phone is duplicate
	if err == nil {
		errMsg := make(map[string]string)
		if len(_u.Name) != 0 && _u.Name == req.Name {
			errMsg["name"] = "name is duplicate"
		}
		if len(_u.Phone) != 0 && _u.Phone == req.Phone {
			errMsg["phone"] = "phone is duplicate"
		}
		if len(_u.Email) != 0 && _u.Email == req.Email {
			errMsg["email"] = "email is duplicate"
		}

		return kerrors.BadRequest("account duplicate", "").WithMetadata(errMsg)
	}

	userInfo := &user.User{
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		PasswordHashed: utils.HashPassword(req.Password, au.conf.PasswordSalt),
		Roles:          []string{req.Role},
	}
	// save
	err = au.userRepo.Save(ctx, userInfo)
	if err != nil {
		return err
	}

	if au.hooks != nil {
		if h, ok := au.hooks[user.ON_REGISTER_EMAIL_VERIFY]; ok {
			if _, err = h(ctx, userInfo); err != nil {
				return err
			}
		}
		if h, ok := au.hooks[user.ON_REGISTER_PHONE_VERIFY]; ok {
			if _, err = h(ctx, userInfo); err != nil {
				return err
			}
		}
		if h, ok := au.hooks[user.ON_REGISTER_SUCCESS]; ok {
			if _, err = h(ctx, userInfo); err != nil {
				return err
			}
		}

	}

	// notify

	return nil
}

func (au *AuthUsecase) Login(ctx context.Context, req *pb.LoginRequest) (*user.User, string, error) {
	if _, err := mail.ParseAddress(req.Name); err == nil {
		req.Email = req.Name
		req.Name = ""
	}

	// find user
	userInfo, err := au.userRepo.Find(ctx, req.Name, req.Email, req.Phone)
	if err != nil {
		return nil, "", err
	}
	// check password
	passwordHashed := utils.HashPassword(req.Password, au.conf.PasswordSalt)
	if userInfo.PasswordHashed != passwordHashed {
		return nil, "", ErrLoginPassword
	}

	if au.hooks != nil {
		if h, ok := au.hooks[user.ON_LOGIN_EMAIL_VERIFY]; ok {
			if !userInfo.EmailVerified {
				return nil, "", ErrEmailNeedVerified
			}
			if h != nil {
				if _, err = h(ctx, userInfo); err != nil {
					return nil, "", err
				}
			}
		}

		if h, ok := au.hooks[user.ON_LOGIN_PHONE_VERIFY]; ok {
			if !userInfo.PhoneVerified {
				return nil, "", ErrPhoneNeedVerified
			}
			if h != nil {
				if _, err = h(ctx, userInfo); err != nil {
					return nil, "", err
				}
			}
		}
	}

	// generate token
	tokenStr, err := utils.GenToken(&utils.UserSession{UserId: userInfo.ID}, au.conf.Jwtsecret)
	if err != nil {
		return nil, "", ErrGenerateToken
	}

	return userInfo, tokenStr, nil
}

func (au *AuthUsecase) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequest) error {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return err
	}
	userInfo, err := au.userRepo.Find(ctx, "", req.Email, req.Phone)
	if err != nil {
		return err
	}

	if au.hooks != nil {
		if h, ok := au.hooks[user.ON_FORGET_PASSWORD]; ok {
			if _, err = h(ctx, userInfo); err != nil {
				return err
			}
		}
	}

	return nil
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) error {
	if req.Code == "" || len(req.Password) < 6 || req.Password != req.RepeatedPassword {
		return ErrLoginPassword
	}

	if au.hooks == nil {
		return ErrNeedHook
	}
	ctx = context.WithValue(ctx, "code", req.Code)

	h, okHooks := au.hooks[user.ON_RESET_PASSWORD]
	if !okHooks {
		return ErrNeedHook
	}
	var err error
	ctx, err = h(ctx, nil)
	if err != nil {
		return err
	}
	userInfo, ok := ctx.Value("resource_user").(*user.User)
	if !ok || userInfo == nil {
		return ErrUserEmpty
	}

	userInfo.PasswordHashed = utils.HashPassword(req.Password, au.conf.PasswordSalt)
	if err := au.userRepo.Save(ctx, userInfo); err != nil {
		return err
	}

	return nil
}

var (
	ErrAccountDuplicate    = errors.New("account duplicate")
	ErrRegisterIsNotPublic = errors.New("register is not public")
	ErrLoginPassword       = errors.New("password error")
	ErrGenerateToken       = errors.New("gen token error")
	ErrEmailNeedVerified   = errors.New("email need verified")
	ErrPhoneNeedVerified   = errors.New("phone need verified")
	ErrNeedHook            = errors.New("need hook")
	ErrUserEmpty           = errors.New("user is empty")
)
