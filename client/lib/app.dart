import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'core/theme/theme.dart';
import 'features/auth/bloc/auth_bloc.dart';
import 'features/auth/bloc/auth_event.dart';
import 'features/auth/bloc/auth_state.dart';
import 'features/auth/ui/login_page.dart';
import 'features/home/home_page.dart';

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: AppTheme.light,
      home: BlocBuilder<AuthBloc, AuthState>(
        builder: (_, state) {
          if (state.status == AuthStatus.authenticated) return const HomePage();
          if (state.status == AuthStatus.initial) {
            context.read<AuthBloc>().add(AuthCheckStatus());
            return const Scaffold(body: Center(child: CircularProgressIndicator()));
          }
          return const LoginPage();
        },
      ),
    );
  }
}
