import 'package:flutter_bloc/flutter_bloc.dart';
import '../data/auth_repository.dart';
import 'auth_event.dart';
import 'auth_state.dart';

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  final AuthRepository _repo;
  AuthBloc(this._repo) : super(const AuthState()) {
    on<AuthLoginRequested>(_onLogin);
    on<AuthRegisterRequested>(_onRegister);
    on<AuthCheckStatus>(_onCheckStatus);
    on<AuthLogoutRequested>(_onLogout);
  }

  Future<void> _onLogin(AuthLoginRequested event, Emitter<AuthState> emit) async {
    emit(const AuthState(status: AuthStatus.loading));
    try {
      await _repo.login(event.phone, event.password);
      emit(const AuthState(status: AuthStatus.authenticated));
    } catch (e) {
      emit(AuthState(status: AuthStatus.error, error: 'login failed'));
    }
  }

  Future<void> _onRegister(AuthRegisterRequested event, Emitter<AuthState> emit) async {
    emit(const AuthState(status: AuthStatus.loading));
    try {
      await _repo.register(event.phone, event.password, event.nickname);
      emit(const AuthState(status: AuthStatus.authenticated));
    } catch (e) {
      emit(AuthState(status: AuthStatus.error, error: 'register failed'));
    }
  }

  Future<void> _onCheckStatus(AuthCheckStatus event, Emitter<AuthState> emit) async {
    final ok = await _repo.isLoggedIn();
    emit(AuthState(status: ok ? AuthStatus.authenticated : AuthStatus.unauthenticated));
  }

  Future<void> _onLogout(AuthLogoutRequested event, Emitter<AuthState> emit) async {
    await _repo.logout();
    emit(const AuthState(status: AuthStatus.unauthenticated));
  }
}
