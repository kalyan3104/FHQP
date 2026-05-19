import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';

class LocationScreen extends StatefulWidget {
  const LocationScreen({super.key});

  @override
  State<LocationScreen> createState() =>
      _LocationScreenState();
}

class _LocationScreenState
    extends State<LocationScreen> {

  String locationMessage =
      "Detecting your location...";

  @override
  void initState() {
    super.initState();

    detectLocation();
  }

  Future<void> detectLocation() async {

    bool serviceEnabled;
    LocationPermission permission;

    // CHECK GPS
    serviceEnabled =
        await Geolocator.isLocationServiceEnabled();

    if (!serviceEnabled) {

      setState(() {
        locationMessage =
            "Location services are disabled.";
      });

      return;
    }

    // CHECK PERMISSION
    permission =
        await Geolocator.checkPermission();

    if (permission == LocationPermission.denied) {

      permission =
          await Geolocator.requestPermission();

      if (permission ==
          LocationPermission.denied) {

        setState(() {
          locationMessage =
              "Location permission denied.";
        });

        return;
      }
    }

    if (permission ==
        LocationPermission.deniedForever) {

      setState(() {
        locationMessage =
            "Location permission permanently denied.";
      });

      return;
    }

    // GET LOCATION
    Position position =
        await Geolocator.getCurrentPosition(
      desiredAccuracy: LocationAccuracy.high,
    );

    setState(() {

      locationMessage =
          "Latitude: ${position.latitude}\n"
          "Longitude: ${position.longitude}";
    });

    // HERE LATER:
    // Navigate to Home Screen
  }

  @override
  Widget build(BuildContext context) {

    return Scaffold(

      backgroundColor: const Color(0xFFF7FAFF),

      body: SafeArea(
        child: Center(

          child: Padding(
            padding: const EdgeInsets.all(24),

            child: Column(
              mainAxisAlignment:
                  MainAxisAlignment.center,

              children: [

                // MAP IMAGE
                Container(
                  width: 260,
                  height: 260,

                  decoration: BoxDecoration(
                    color: Colors.white,

                    borderRadius:
                        BorderRadius.circular(24),

                    boxShadow: [
                      BoxShadow(
                        color:
                            Colors.black.withOpacity(0.05),

                        blurRadius: 10,
                      ),
                    ],
                  ),

                  child: Center(
                    child: Icon(
                      Icons.location_on,
                      size: 120,
                      color: Color(0xFF12003B),
                    ),
                  ),
                ),

                const SizedBox(height: 30),

                Text(
                  locationMessage,

                  textAlign: TextAlign.center,

                  style: const TextStyle(
                    fontSize: 18,
                    color: Colors.black87,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}