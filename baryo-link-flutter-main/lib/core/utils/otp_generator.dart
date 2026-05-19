import 'dart:math';

class OtpGenerator {

  static String generateOTP() {

    Random random = Random();

    int otp = 1000 + random.nextInt(9000);

    return otp.toString();
  }
}