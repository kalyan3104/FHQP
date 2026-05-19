import 'dart:ui';
import 'package:flutter/material.dart';

class GlassStyles {
  static const double blurSigma = 16.0;
  static const double borderRadius = 20.0;

  // Adaptive decoration based on context
  static BoxDecoration glassDecoration(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    return BoxDecoration(
      color: isDark 
          ? const Color(0xFF1E293B).withOpacity(0.7) 
          : Colors.white.withOpacity(0.6), // White glass for light mode
      borderRadius: BorderRadius.circular(borderRadius),
      border: Border.all(
        color: isDark ? Colors.white24 : Colors.white, 
        width: 1.5,
      ),
      boxShadow: [
        BoxShadow(
          color: isDark ? Colors.black.withOpacity(0.1) : Colors.indigo.withOpacity(0.05),
          blurRadius: 20,
          spreadRadius: 2,
        ),
      ],
    );
  }

  static Widget glassContainer({
    required Widget child,
    required BuildContext context,
    double? radius,
    EdgeInsetsGeometry? padding,
    double? width,
    double? height,
  }) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(radius ?? borderRadius),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: blurSigma, sigmaY: blurSigma),
        child: Container(
          width: width,
          height: height,
          padding: padding,
          decoration: glassDecoration(context).copyWith(
            borderRadius: BorderRadius.circular(radius ?? borderRadius),
          ),
          child: child,
        ),
      ),
    );
  }
}

class GlassCard extends StatelessWidget {
  final Widget child;
  final EdgeInsetsGeometry? padding;
  final double? width;
  final double? height;
  final VoidCallback? onTap;

  const GlassCard({
    super.key,
    required this.child,
    this.padding,
    this.width,
    this.height,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    Widget content = Container(
      width: width,
      height: height,
      padding: padding ?? const EdgeInsets.all(16),
      decoration: GlassStyles.glassDecoration(context),
      child: child,
    );

    if (onTap != null) {
      content = InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(GlassStyles.borderRadius),
        child: content,
      );
    }

    return ClipRRect(
      borderRadius: BorderRadius.circular(GlassStyles.borderRadius),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: GlassStyles.blurSigma, sigmaY: GlassStyles.blurSigma),
        child: content,
      ),
    );
  }
}
