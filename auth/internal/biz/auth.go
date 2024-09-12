package biz

import (
	"context"
	"errors"
	pb "github.com/databonfire/bonfire/auth/api/v1"
	"github.com/databonfire/bonfire/auth/internal/conf"
	"github.com/databonfire/bonfire/auth/internal/utils"
	"github.com/databonfire/bonfire/auth/user"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"net/mail"
	"time"
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
	errMsg := make(map[string]string)
	//if len(req.Email) == 0 {
	//	// Phone and email cannot both be empty
	//	errMsg = map[string]string{
	//		"email": "The email format is incorrect",
	//		//"phone": "phone is empty",
	//	}
	//	return kerrors.BadRequest("phone and email cannot both be empty", "").WithMetadata(errMsg)
	//}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		errMsg["email"] = "The email format is incorrect"
	}
	if len(req.Name) < 3 || len(req.Name) > 30 {
		errMsg["name"] = "Keep the length between 3 and 30 characters"
	}
	if len(req.Password) < 6 || len(req.Password) > 12 {
		errMsg["password"] = "Password length should be between 6 and 12 characters"
	}
	if req.Password != req.Repassword {
		errMsg["password"] = "Passwords do not match"
		errMsg["repassword"] = "Passwords do not match"
	}
	if len(req.CompanyName) < 3 || len(req.CompanyName) > 50 {
		errMsg["company_name"] = "Keep the length between 3 and 50 characters"
	}

	if len(errMsg) != 0 {
		return kerrors.BadRequest(ReasonInvalidParam, "").WithMetadata(errMsg)
	}

	_u, err := au.userRepo.Find(ctx, req.Name, req.Email, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// name, email, phone is duplicate
	if err == nil {

		if len(_u.Name) != 0 && _u.Name == req.Name {
			errMsg["name"] = "Username has already been taken"
		}
		if len(_u.Phone) != 0 && _u.Phone == req.Phone {
			errMsg["phone"] = "Phone has already been taken"
		}
		if len(_u.Email) != 0 && _u.Email == req.Email {
			errMsg["email"] = "This email has already been registered"
		}

		return kerrors.BadRequest(ReasonInvalidParam, "").WithMetadata(errMsg)
		//if _u.EmailVerified {
		//	return kerrors.BadRequest(ReasonInvalidParam, "").WithMetadata(errMsg)
		//} else {
		//	return kerrors.BadRequest(ReasonNeedVerify, "").WithMetadata(errMsg)
		//}
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
	errMsg := make(map[string]string)
	if _, err := mail.ParseAddress(req.Name); err == nil {
		req.Email = req.Name
		req.Name = ""
	}

	// find user
	userInfo, err := au.userRepo.Find(ctx, req.Name, req.Email, req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errMsg["name"] = "Account does not exist, please sign up"
			return nil, "", kerrors.BadRequest(ReasonInvalidParam, "Account does not exist, please sign up").WithMetadata(errMsg)
		}
		return nil, "", pb.ErrorInternal(err.Error())
	}
	// check password
	passwordHashed := utils.HashPassword(req.Password, au.conf.PasswordSalt)
	if userInfo.PasswordHashed != passwordHashed {
		errMsg["password"] = "Incorrect password"
		return nil, "", kerrors.BadRequest(ReasonInvalidParam, "Incorrect password").WithMetadata(errMsg)
	}
	if userInfo.ExpiredAt <= time.Now().Unix() {
		//return nil, "", kerrors.BadRequest(ReasonPaymentRequired, "You account is expired,please email service@kolplanet.com to renew your package").WithMetadata(errMsg)
		return nil, "", kerrors.BadRequest(ReasonPaymentRequired, "KOLPlanet team will contact you soon.").WithMetadata(errMsg)
	}

	if au.hooks != nil {
		if h, ok := au.hooks[user.ON_LOGIN_EMAIL_VERIFY]; ok {
			if !userInfo.EmailVerified {
				return nil, "", kerrors.BadRequest(ReasonNeedVerify, "").WithMetadata(errMsg)
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

func (au *AuthUsecase) ResendOTP(ctx context.Context, req *pb.ResendOTPRequest) error {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return err
	}
	userInfo, err := au.userRepo.Find(ctx, "", req.Email, req.Phone)
	if err != nil {
		errMsg := make(map[string]string)
		errMsg["email"] = "Account does not exist, please sign up"
		return kerrors.BadRequest(ReasonInvalidParam, "Resent register error").WithMetadata(errMsg)
		//return err
	}

	if au.hooks != nil {
		if h, ok := au.hooks[user.ON_RESEND_OTP]; ok {
			if _, err = h(ctx, userInfo); err != nil {
				return err
			}
		}
	}

	return nil
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) error {
	if req.Code == "" {
		return ErrResetPasswordCode
	}
	if len(req.Password) < 6 || len(req.Password) > 12 || req.Password != req.RepeatedPassword {
		return ErrResetPassword
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
	ErrAccountNotExist     = pb.ErrorInvalidParam("Account does not exist, please sign up")
	ErrLoginPassword       = pb.ErrorInvalidParam("Incorrect password")
	ErrResetPasswordCode   = pb.ErrorInvalidParam("Code is error")
	ErrResetPassword       = pb.ErrorInvalidParam("Password length should be between 6 and 12 characters.Or passwords do not match")
	ErrGenerateToken       = pb.ErrorInternal("gen token error")
	ErrEmailNeedVerified   = pb.ErrorInvalidParam("email need verified")
	ErrPhoneNeedVerified   = pb.ErrorInvalidParam("phone need verified")
	ErrNeedHook            = pb.ErrorInvalidParam("need hook")
	ErrUserEmpty           = pb.ErrorInvalidParam("user is empty")
)

var (
	ErrNoEnoughCredit = pb.ErrorPaymentRequired("no enough credits")
	ErrExternal       = pb.ErrorExternal("external err")
	ErrInternal       = pb.ErrorInternal("internal err")
	ErrParam          = pb.ErrorInvalidParam("external err")
)

const (
	ReasonInvalidParam     = "INVALID_PARAM"
	ReasonNeedLogin        = "NEED_LOGIN"
	ReasonPaymentRequired  = "PAYMENT_REQUIRED"
	ReasonForbiddenRequest = "FORBIDDEN_REQUEST"
	ReasonInternal         = "INTERNAL"
	ReasonExternal         = "EXTERNAL"
	ReasonNeedVerify       = "NEED_VERIFY"
)
