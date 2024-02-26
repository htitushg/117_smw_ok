package controllers

//import "118_session_ok/controllers"

var HomeBundle = Join(Home, Log())
var LoginBundle = Join(Login, Log())
var SigninBundle = Join(Signin, Log())
var IndexBundle = Join(Index, Log())
var LogoutBundle = Join(Logout, Log(), Foo())
var RegisterBundle = Join(Register, Log(), Foo())
var AfficheUserInfoBundle = Join(AfficheUserInfo, Log(), Foo())
var IndexHandlerNoMethBundle = Join(IndexHandlerNoMeth, Log(), Foo())
var IndexHandlerOtherBundle = Join(IndexHandlerOther, Log(), Foo())
