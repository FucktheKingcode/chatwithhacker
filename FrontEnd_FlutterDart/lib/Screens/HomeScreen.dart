import 'package:chat_app/Authenticate/Methods.dart';
import 'package:chat_app/Screens/GiaiMaScreen.dart';
import 'package:chat_app/Screens/MaHoaScreen.dart';
import 'package:chat_app/Screens/SearchScreen.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  @override
  _HomeScreenState createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('My Screen'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            ElevatedButton(
              onPressed: () => Navigator.of(context)
                  .push(MaterialPageRoute(builder: (_) => SearchScreen())),
              child: Text('Chat'),
              style: ButtonStyle(
                padding: MaterialStateProperty.all<EdgeInsets>(
                  EdgeInsets.symmetric(
                      horizontal: 158,
                      vertical: 50), // Điều chỉnh kích thước của nút
                ),
              ),
            ),
            SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => Navigator.of(context)
                  .push(MaterialPageRoute(builder: (_) => MaHoaScreen())),
              child: Text('Mã hóa'),
              style: ButtonStyle(
                padding: MaterialStateProperty.all<EdgeInsets>(
                  EdgeInsets.symmetric(
                      horizontal: 150,
                      vertical: 50), // Điều chỉnh kích thước của nút
                ),
              ),
            ),
            SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => Navigator.of(context)
                  .push(MaterialPageRoute(builder: (_) => GiaiMaScreen())),
              child: Text('Giải Mã'),
              style: ButtonStyle(
                padding: MaterialStateProperty.all<EdgeInsets>(
                  EdgeInsets.symmetric(
                      horizontal: 150,
                      vertical: 50), // Điều chỉnh kích thước của nút
                ),
              ),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => logOut(context),
        child: Icon(Icons.logout),
      ),
    );
  }
}
