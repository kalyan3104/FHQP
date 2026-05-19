import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../features/auth/presentation/screens/login_screen.dart';
import '../../features/auth/presentation/screens/otp_screen.dart';
import '../../features/home/presentation/screens/home_screen.dart';
import '../../features/home/presentation/screens/main_screen.dart';
import '../../features/booking/presentation/screens/booking_screen.dart';

// Placeholder screens for tabs
class PlaceholderScreen extends StatelessWidget {
  final String title;
  const PlaceholderScreen(this.title, {super.key});
  @override 
  Widget build(BuildContext context) => Scaffold(body: Center(child: Text(title)));
}

final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/login',
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/otp',
        builder: (context, state) => const OtpScreen(),
      ),
      ShellRoute(
        navigatorKey: _shellNavigatorKey,
        builder: (context, state, child) => MainScreen(child: child),
        routes: [
          GoRoute(
            path: '/home',
            builder: (context, state) => const HomeScreen(),
          ),
          GoRoute(
            path: '/bookings',
            builder: (context, state) => const PlaceholderScreen("My Bookings"),
          ),
          GoRoute(
            path: '/calendar',
            builder: (context, state) => const PlaceholderScreen("Calendar"),
          ),
          GoRoute(
            path: '/account',
            builder: (context, state) => const PlaceholderScreen("Account"),
          ),
        ],
      ),
      GoRoute(
        path: '/booking',
        parentNavigatorKey: _rootNavigatorKey, // Push over Shell
        builder: (context, state) => const BookingScreen(),
      ),
    ],
  );
});
