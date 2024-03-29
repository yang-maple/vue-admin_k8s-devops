package service

import (
	"kubeops/dao"
	"kubeops/model"
	"kubeops/utils"
)

var Login login

type login struct {
}

type LoginInfo struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	CaptchaId   string `json:"captcha_id"`
	VerifyValue string `json:"verify_value"`
}

type UserInfo struct {
	Id          int              `json:"id"`
	UserName    string           `json:"user_name"`
	Email       string           `json:"email"`
	Roles       []model.UserRole `json:"roles"`
	Avatar      string           `json:"avatar"`
	ClusterName string           `json:"cluster_name"`
}

// VerifyUserInfo 验证用户身份信息
func (l *login) VerifyUserInfo(info *LoginInfo) (token string, err error) {
	uuid, err := dao.Login.VerifyUser(info.Username, info.Password)
	if err != nil {
		utils.Logger.Error("User identity verification failed")
		return "", err
	}
	utils.Logger.Info("User identity verification succeeded")
	token, err = utils.CreateJwtToken(uuid, info.Username, utils.UserExpireDuration, utils.UserSecret)
	if err != nil {
		utils.Logger.Error("CreateJwtToken failed :" + err.Error())
		return "", err
	}
	return token, nil
}

// GetUserInfo 获取用户信息
func (l *login) GetUserInfo(token string) (info *UserInfo, err error) {
	//解析token
	claim, err := utils.JWTToken.ParseToken(token, utils.UserSecret)
	if err != nil {
		utils.Logger.Error("ParseJwtToken failed :" + err.Error())
		return nil, err
	}
	//获取用户信息
	userinfo, err := dao.Login.GetUserInfo(claim.Id)
	if err != nil {
		utils.Logger.Error("GetUserInfo failed :" + err.Error())
		return nil, err
	}
	//获取用户集群信息
	url, cname, err := dao.Login.GetUserCluster(claim.Id)
	if err != nil {
		utils.Logger.Error("GetUserCluster failed :" + err.Error())
		return nil, err
	}
	// 如果客户端已存在则不初始化
	if K8s.ConfigDir[claim.Id] != nil {
		return &UserInfo{
			Id:          int(userinfo.Id),
			UserName:    userinfo.Username,
			Avatar:      userinfo.Avatar,
			Roles:       []model.UserRole{userinfo.Roles},
			ClusterName: cname,
		}, nil
	}
	//如果存在集群则初始化
	if cname != "" {
		//存入集群信息
		K8s.ConfigDir[claim.Id] = &url
		//初始化
		err = K8s.Init(claim.Id)
		if err != nil {
			return nil, err
		}
	}
	//返回用户信息
	return &UserInfo{
		Id:          int(userinfo.Id),
		UserName:    userinfo.Username,
		Avatar:      userinfo.Avatar,
		Roles:       []model.UserRole{userinfo.Roles},
		ClusterName: cname,
	}, nil
}
