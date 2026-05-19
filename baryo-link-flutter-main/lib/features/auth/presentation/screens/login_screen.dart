import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:country_code_picker/country_code_picker.dart';
import 'package:gap/gap.dart';
import 'package:phone_numbers_parser/phone_numbers_parser.dart';

import 'otp_screen.dart';
import '../../../../core/utils/otp_generator.dart';
import '../../../../core/theme/glass_styles.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() =>
      _LoginScreenState();
}

class _LoginScreenState
    extends ConsumerState<LoginScreen> {

  final TextEditingController _phoneController =
      TextEditingController();

  String _selectedCountryCode = 'US';

  @override
  void dispose() {
    _phoneController.dispose();
    super.dispose();
  }

  // PHONE VALIDATION
  bool isValidPhone(
    String phone,
    String isoCode,
  ) {
    try {

      final parsed = PhoneNumber.parse(
        phone,
        destinationCountry:
            IsoCode.values.firstWhere(
          (e) => e.name == isoCode,
        ),
      );

      return parsed.isValid();

    } catch (e) {

      return false;
    }
  }

  // LOGIN
  void _onLogin() {

    final phone =
        _phoneController.text.trim();

    final isValid = isValidPhone(
      phone,
      _selectedCountryCode,
    );

    // INVALID
    if (!isValid) {

      ScaffoldMessenger.of(context)
          .showSnackBar(
        const SnackBar(
          content: Text(
            'Enter valid phone number',
          ),
          backgroundColor: Colors.red,
        ),
      );

      return;
    }

    // GENERATE OTP
    final otp =
        OtpGenerator.generateOTP();

    print('Generated OTP: $otp');

    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (_) => OtpScreen(
          otp: otp,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(

      backgroundColor:
          const Color(0xFFF8FAFF),

      body: SafeArea(

        child: SingleChildScrollView(

          child: Padding(
            padding:
                const EdgeInsets.all(24),

            child: Column(
              crossAxisAlignment:
                  CrossAxisAlignment.start,

              children: [

                // TOP RIGHT
                const Align(
                  alignment:
                      Alignment.topRight,

                  child: Text(
                    'Explore app',

                    style: TextStyle(
                      fontWeight:
                          FontWeight.w600,

                      color:
                          Color(0xFF1E1B4B),
                    ),
                  ),
                ),

                const SizedBox(
                  height: 70,
                ),

                // TITLE
                const Text(
                  "Let's get you started",

                  style: TextStyle(
                    fontSize: 38,
                    fontWeight:
                        FontWeight.bold,

                    color:
                        Color(0xFF1E1B4B),
                  ),
                ),

                const SizedBox(
                  height: 16,
                ),

                // SUBTITLE
                Text(
                  'Enter your mobile number',

                  style: TextStyle(
                    fontSize: 18,
                    fontWeight:
                        FontWeight.w500,

                    color:
                        Colors.grey[700],
                  ),
                ),

                const SizedBox(
                  height: 12,
                ),

                Text(
                  'We\'ll send you a one-time code to verify your mobile number.',

                  style: TextStyle(
                    fontSize: 15,
                    height: 1.5,
                    color:
                        Colors.grey[500],
                  ),
                ),

                const SizedBox(
                  height: 40,
                ),

                // PHONE CARD
                GlassCard(
                  padding:
                      const EdgeInsets.symmetric(
                    horizontal: 16,
                    vertical: 12,
                  ),

                  child: Row(
                    children: [

                      // COUNTRY PICKER
                      Container(
                        decoration:
                            BoxDecoration(
                          color:
                              Colors.white,

                          borderRadius:
                              BorderRadius.circular(
                                  12),

                          border: Border.all(
                            color: Colors
                                .grey.shade300,
                          ),
                        ),

                        child:
                            CountryCodePicker(

                          onChanged:
                              (country) {

                            setState(() {

                              _selectedCountryCode =
                                  country.code ??
                                      'US';
                            });
                          },

                          initialSelection:
                              'US',

                          favorite:
                              const [
                            '+1',
                            'IN',
                          ],

                          showCountryOnly:
                              false,

                          showOnlyCountryWhenClosed:
                              false,

                          alignLeft:
                              false,
                        ),
                      ),

                      const SizedBox(
                        width: 16,
                      ),

                      Container(
                        width: 1,
                        height: 24,
                        color: Colors
                            .grey.shade300,
                      ),

                      const SizedBox(
                        width: 16,
                      ),

                      // TEXTFIELD
                      Expanded(
                        child: TextField(

                          controller:
                              _phoneController,

                          keyboardType:
                              TextInputType.phone,

                          decoration:
                              const InputDecoration(

                            hintText:
                                'Mobile number',

                            border:
                                InputBorder.none,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),

                const SizedBox(
                  height: 260,
                ),

                // TERMS
                Text(
                  'By continuing, you agree to ServU\'s Terms of Use and Privacy Policy.',

                  textAlign:
                      TextAlign.center,

                  style: TextStyle(
                    fontSize: 12,
                    color:
                        Colors.grey[500],
                  ),
                ),

                const SizedBox(
                  height: 20,
                ),

                // BUTTON
                SizedBox(
                  width: double.infinity,
                  height: 56,

                  child: ElevatedButton(

                    onPressed:
                        _onLogin,

                    style:
                        ElevatedButton
                            .styleFrom(
                      backgroundColor:
                          const Color(
                              0xFF1E1B4B),

                      shape:
                          RoundedRectangleBorder(
                        borderRadius:
                            BorderRadius.circular(
                                14),
                      ),
                    ),

                    child: const Text(
                      'Continue',

                      style: TextStyle(
                        fontSize: 16,
                        fontWeight:
                            FontWeight.bold,

                        color:
                            Colors.white,
                      ),
                    ),
                  ),
                ),

                const SizedBox(
                  height: 20,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}