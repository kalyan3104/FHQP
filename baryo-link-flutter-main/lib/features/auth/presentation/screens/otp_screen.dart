import 'package:flutter/material.dart';
import 'package:gap/gap.dart';
import 'package:pinput/pinput.dart';

import 'location_screen.dart';

class OtpScreen extends StatefulWidget {
  final String otp;

  const OtpScreen({
    super.key,
    required this.otp,
  });

  @override
  State<OtpScreen> createState() => _OtpScreenState();
}

class _OtpScreenState extends State<OtpScreen> {

  final TextEditingController _otpController =
      TextEditingController();

  void verifyOtp() {

    if (_otpController.text == widget.otp) {

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('OTP Verified'),
          backgroundColor: Colors.green,
        ),
      );

      Navigator.pushReplacement(
        context,
        MaterialPageRoute(
          builder: (_) => const LocationScreen(),
        ),
      );

    } else {

      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Invalid OTP'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(

      backgroundColor: const Color(0xFFF7FAFF),

      body: SafeArea(

        child: SingleChildScrollView(

          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: 24,
              vertical: 20,
            ),

            child: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,

              children: [

                // TOP BAR
                Row(
                  mainAxisAlignment:
                      MainAxisAlignment.spaceBetween,

                  children: [

                    IconButton(
                      onPressed: () {
                        Navigator.pop(context);
                      },

                      icon: const Icon(
                        Icons.arrow_back,
                        color: Color(0xFF12003B),
                      ),
                    ),

                    const Text(
                      'Explore app',

                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: Color(0xFF12003B),
                      ),
                    ),
                  ],
                ),

                const Gap(35),

                // TITLE
                const Text(
                  'OTP\nVerification',

                  style: TextStyle(
                    fontSize: 38,
                    height: 1.2,
                    fontWeight: FontWeight.bold,
                    color: Color(0xFF12003B),
                  ),
                ),

                const Gap(16),

                // SUBTITLE
                const Text(
                  'Enter the OTP sent to your\nmobile number',

                  style: TextStyle(
                    color: Colors.grey,
                    fontSize: 16,
                    height: 1.5,
                  ),
                ),

                const Gap(40),

                // OTP INPUT
                Center(
                  child: Pinput(

                    controller: _otpController,

                    length: 6,

                    mainAxisAlignment:
                        MainAxisAlignment.spaceBetween,

                    defaultPinTheme: PinTheme(
                      width: 48,
                      height: 58,

                      textStyle: const TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: Color(0xFF12003B),
                      ),

                      decoration: BoxDecoration(
                        color: Colors.white,

                        borderRadius:
                            BorderRadius.circular(12),

                        border: Border.all(
                          color: Colors.grey.shade300,
                        ),
                      ),
                    ),
                  ),
                ),

                const Gap(35),

                // TIMER
                const Center(
                  child: Text(
                    '00:30',

                    style: TextStyle(
                      fontWeight: FontWeight.bold,
                      color: Color(0xFF12003B),
                    ),
                  ),
                ),

                const Gap(14),

                // RESEND
                Center(
                  child: RichText(
                    text: const TextSpan(

                      text: "Didn't receive OTP? ",

                      style: TextStyle(
                        color: Colors.grey,
                        fontSize: 14,
                      ),

                      children: [

                        TextSpan(
                          text: 'Resend OTP',

                          style: TextStyle(
                            color: Color(0xFF12003B),
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),

                const Gap(60),

                // CONTINUE BUTTON
                SizedBox(
                  width: double.infinity,
                  height: 56,

                  child: ElevatedButton(

                    onPressed: verifyOtp,

                    style: ElevatedButton.styleFrom(
                      backgroundColor:
                          const Color(0xFF12003B),

                      shape: RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius.circular(14),
                      ),
                    ),

                    child: const Text(
                      'Continue',

                      style: TextStyle(
                        fontSize: 16,
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),

                const Gap(25),
              ],
            ),
          ),
        ),
      ),
    );
  }
}