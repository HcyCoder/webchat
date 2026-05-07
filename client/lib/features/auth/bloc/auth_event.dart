abstract class AuthEvent {}
class AuthLoginRequested extends AuthEvent {
  final String phone, password;
  AuthLoginRequested(this.phone, this.password);
}
class AuthRegisterRequested extends AuthEvent {
  final String phone, password, nickname;
  AuthRegisterRequested(this.phone, this.password, this.nickname);
}
class AuthCheckStatus extends AuthEvent {}
class AuthLogoutRequested extends AuthEvent {}
