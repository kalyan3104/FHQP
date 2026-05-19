import 'dart:async';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {

  final List<Map<String, String>> services = [

    {
      'image': 'assets/images/auto_repair.png',
      'title': 'Auto repair\nservice',
    },

    {
      'image': 'assets/images/home_cleaning.png',
      'title': 'Home cleaning\nservices',
    },

    {
      'image': 'assets/images/hvac.png',
      'title': 'HVAC\nMaintenance',
    },

    {
      'image': 'assets/images/lawn.png',
      'title': 'Lawn Care &\nGardening',
    },

    {
      'image': 'assets/images/moving.png',
      'title': 'Moving\nAssistance',
    },

    {
      'image': 'assets/images/electrical.png',
      'title': 'Electrical\nServices',
    },
  ];

  final List<String> banners = [
    'assets/images/banner.png',
    'assets/images/banner2.png',
    'assets/images/banner3.png',
  ];

  final PageController _pageController = PageController();

  int currentPage = 0;

  Timer? timer;

  @override
  void initState() {
    super.initState();

    timer = Timer.periodic(
      const Duration(seconds: 3),
      (timer) {

        if (currentPage < banners.length - 1) {
          currentPage++;
        } else {
          currentPage = 0;
        }

        _pageController.animateToPage(
          currentPage,
          duration: const Duration(milliseconds: 400),
          curve: Curves.easeInOut,
        );

        setState(() {});
      },
    );
  }

  @override
  void dispose() {
    timer?.cancel();
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(

      backgroundColor: const Color(0xFFF5F7FB),

      bottomNavigationBar: BottomNavigationBar(

        currentIndex: 0,

        selectedItemColor: const Color(0xFF12003B),

        unselectedItemColor: Colors.grey,

        type: BottomNavigationBarType.fixed,

        items: const [

          BottomNavigationBarItem(
            icon: Icon(Icons.home_outlined),
            label: 'Home',
          ),

          BottomNavigationBarItem(
            icon: Icon(Icons.grid_view_outlined),
            label: 'Categories',
          ),

          BottomNavigationBarItem(
            icon: Icon(Icons.calendar_month_outlined),
            label: 'Bookings',
          ),

          BottomNavigationBarItem(
            icon: Icon(Icons.person_outline),
            label: 'Profile',
          ),
        ],
      ),

      body: SafeArea(

        child: SingleChildScrollView(

          child: Padding(

            padding: const EdgeInsets.all(16),

            child: Column(

              crossAxisAlignment: CrossAxisAlignment.start,

              children: [

                /// HEADER
                Row(

                  mainAxisAlignment:
                      MainAxisAlignment.spaceBetween,

                  children: [

                    Column(

                      crossAxisAlignment:
                          CrossAxisAlignment.start,

                      children: const [

                        Text(
                          'ServU',

                          style: TextStyle(
                            fontSize: 28,
                            fontWeight: FontWeight.bold,
                            color: Color(0xFF12003B),
                          ),
                        ),

                        SizedBox(height: 4),

                        Row(
                          children: [

                            Icon(
                              Icons.location_on_outlined,
                              size: 16,
                              color: Colors.grey,
                            ),

                            SizedBox(width: 4),

                            Text(
                              'ABC Residency...',
                              style: TextStyle(
                                color: Colors.grey,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),

                    const Icon(
                      Icons.notifications_none,
                      size: 28,
                    ),
                  ],
                ),

                const SizedBox(height: 20),

                /// SEARCH BAR
                Container(

                  padding:
                      const EdgeInsets.symmetric(
                    horizontal: 14,
                  ),

                  height: 52,

                  decoration: BoxDecoration(
                    color: Colors.white,

                    borderRadius:
                        BorderRadius.circular(16),
                  ),

                  child: const Row(
                    children: [

                      Icon(
                        Icons.search,
                        color: Colors.grey,
                      ),

                      SizedBox(width: 10),

                      Expanded(
                        child: Text(
                          'Search for services...',
                          style: TextStyle(
                            color: Colors.grey,
                          ),
                        ),
                      ),

                      Icon(
                        Icons.mic_none,
                        color: Colors.grey,
                      ),
                    ],
                  ),
                ),

                const SizedBox(height: 24),

                /// SERVICES
                SizedBox(

                  height: 130,

                  child: ListView.builder(

                    scrollDirection: Axis.horizontal,

                    itemCount: services.length,

                    itemBuilder: (context, index) {

                      final service = services[index];

                      return Padding(

                        padding:
                            const EdgeInsets.only(
                          right: 12,
                        ),

                        child: Container(

                          width: 120,

                          padding:
                              const EdgeInsets.all(12),

                          decoration: BoxDecoration(
                            color: Colors.white,

                            borderRadius:
                                BorderRadius.circular(18),

                            border: Border.all(
                              color: Colors.grey.shade200,
                            ),
                          ),

                          child: Column(

                            mainAxisAlignment:
                                MainAxisAlignment.center,

                            children: [

                              Image.asset(
                                service['image']!,

                                height: 50,

                                fit: BoxFit.contain,
                              ),

                              const SizedBox(height: 10),

                              Text(
                                service['title']!,

                                textAlign: TextAlign.center,

                                style: const TextStyle(
                                  fontSize: 13,
                                  fontWeight:
                                      FontWeight.w500,

                                  color:
                                      Color(0xFF12003B),
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  ),
                ),

                const SizedBox(height: 24),

                /// BANNER SLIDER
Column(
  children: [

    SizedBox(
      height: 158,

      child: PageView.builder(

        controller: _pageController,

        itemCount: banners.length,

        onPageChanged: (index) {

          setState(() {
            currentPage = index;
          });
        },

        itemBuilder: (context, index) {

          return Padding(

            padding: const EdgeInsets.symmetric(
              horizontal: 16,
            ),

            child: Container(

              width: MediaQuery.of(context).size.width,

              decoration: BoxDecoration(

                borderRadius:
                    BorderRadius.circular(22),

                image: DecorationImage(

                  image: AssetImage(
                    banners[index],
                  ),

                  fit: BoxFit.cover,
                ),
              ),
            ),
          );
        },
      ),
    ),

                    const SizedBox(height: 10),

                    Row(

                      mainAxisAlignment:
                          MainAxisAlignment.center,

                      children: List.generate(

                        banners.length,

                        (index) {

                          return AnimatedContainer(

                            duration:
                                const Duration(
                              milliseconds: 300,
                            ),

                            margin:
                                const EdgeInsets.symmetric(
                              horizontal: 4,
                            ),

                            width:
                                currentPage == index
                                    ? 20
                                    : 8,

                            height: 6,

                            decoration: BoxDecoration(

                              color:
                                  currentPage == index
                                      ? const Color(
                                          0xFF12003B)
                                      : Colors
                                          .grey.shade300,

                              borderRadius:
                                  BorderRadius.circular(
                                20,
                              ),
                            ),
                          );
                        },
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}