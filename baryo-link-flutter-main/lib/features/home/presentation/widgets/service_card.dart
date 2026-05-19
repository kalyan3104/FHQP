import 'package:flutter/material.dart';
import 'package:gap/gap.dart';
import '../../../../core/theme/glass_styles.dart';

class ServiceCard extends StatelessWidget {
  final String title;
  final String category;
  final String rating;
  final String price;
  final VoidCallback onTap;

  const ServiceCard({
    super.key,
    required this.title,
    required this.category,
    required this.rating,
    required this.price,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: GlassCard(
        onTap: onTap,
        padding: EdgeInsets.zero,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Placeholder Image Area
            Container(
              height: 120,
              width: double.infinity,
              color: Colors.grey.withOpacity(0.1),
              child: const Center(
                child: Icon(Icons.cleaning_services, size: 48, color: Colors.black12),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        category.toUpperCase(),
                        style: TextStyle(
                          color: Theme.of(context).primaryColor, 
                          fontSize: 10, 
                          fontWeight: FontWeight.bold,
                          letterSpacing: 1.2,
                        ),
                      ),
                      Row(
                        children: [
                          const Icon(Icons.star, size: 14, color: Colors.amber),
                          const Gap(4),
                          Text(rating, style: TextStyle(color: Colors.grey[600], fontSize: 12)),
                        ],
                      ),
                    ],
                  ),
                  const Gap(8),
                  Text(
                    title,
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: const Color(0xFF1E1B4B), // Navy
                    ),
                  ),
                  const Gap(4),
                  Text(
                    price,
                    style: Theme.of(context).textTheme.titleSmall?.copyWith(
                      color: Colors.grey[700],
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
