import 'package:flutter/material.dart';
import 'login_screen.dart';

class WelcomeScreen extends StatelessWidget {
  const WelcomeScreen({super.key});

  @override
  Widget build(BuildContext context) {

    return Scaffold(
      backgroundColor: Colors.white,

      body: SafeArea(

        child: SingleChildScrollView(

          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: 24,
            ),

            child: Column(
              children: [

                const SizedBox(height: 20),

                // IMAGE
                Image.asset(
                  'assets/images/welcome_banner.png',
                  height: 260,
                  fit: BoxFit.contain,
                ),

                const SizedBox(height: 30),

                // TITLE
                const Text(
                  'Welcome to ServU',
                  textAlign: TextAlign.center,

                  style: TextStyle(
                    fontSize: 28,
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF12003B),
                  ),
                ),

                const SizedBox(height: 12),

                // SUBTITLE
                const Text(
                  'Your one-stop solution for finding trusted services.',

                  textAlign: TextAlign.center,

                  style: TextStyle(
                    fontSize: 15,
                    color: Colors.grey,
                    height: 1.5,
                  ),
                ),

                const SizedBox(height: 50),

                // BUTTON
                SizedBox(
                  width: double.infinity,
                  height: 56,

                  child: ElevatedButton(

                    onPressed: () {

                      Navigator.push(
                        context,

                        MaterialPageRoute(
                          builder: (_) =>
                              const LoginScreen(),
                        ),
                      );
                    },

                    style:
                        ElevatedButton.styleFrom(
                      backgroundColor:
                          const Color(0xFF12003B),

                      shape:
                          RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius.circular(14),
                      ),
                    ),

                    child: const Text(
                      'Get started',

                      style: TextStyle(
                        fontSize: 16,
                        color: Colors.white,
                        fontWeight:
                            FontWeight.w600,
                      ),
                    ),
                  ),
                ),

                const SizedBox(height: 18),

                // EXPLORE APP
                const Text(
                  'Explore app',

                  style: TextStyle(
                    fontSize: 15,
                    color: Color(0xFF12003B),
                    fontWeight: FontWeight.w500,
                  ),
                ),

                const SizedBox(height: 30),
              ],
            ),
          ),
        ),
      ),
    );
  }
}