import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class NewPage extends StatefulWidget {
  @override
  _NewPageState createState() => _NewPageState();
}

class _NewPageState extends State<NewPage> {
  final TextEditingController textEditingController = TextEditingController();
  String hash = '';
  String aesM = '';
  String privateKey = '';
  String cipherKey = '';
  String decryptedData = '';
  bool checkData = false;

  @override
  void dispose() {
    // Clean up the controller when the widget is disposed.
    textEditingController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('New Page'),
      ),
      body: SingleChildScrollView(
        padding: EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              'Nơi bỏ vào dữ liệu mã hóa',
              style: TextStyle(fontSize: 20),
            ),
            SizedBox(height: 10),
            TextField(
              controller: textEditingController,
              maxLines: null, // Cho phép nhiều dòng
              keyboardType: TextInputType.multiline, // Hiển thị bàn phím dạng đa dòng
              decoration: InputDecoration(
                hintText: 'Nhập dữ liệu',
                border: OutlineInputBorder(),
              ),
            ),
            SizedBox(height: 20),
            ElevatedButton(
              onPressed: () {
                String inputData = textEditingController.text;
                List<String> dataArray = splitData(inputData);
                if (dataArray.length >= 9) { // Sửa đây, bạn đã chỉ định dataArray.length >= 1
                  setState(() {
                    hash = dataArray[5];
                    aesM = dataArray[6];
                    privateKey = dataArray[8];
                    cipherKey = dataArray[9];
                  });
                  sendDecryptRequest(hash, aesM, privateKey, cipherKey);
                } else {
                  // Xử lý khi dữ liệu không đủ
                }
              },
              child: Text('Giải mã'),
            ),
            SizedBox(height: 20),
            Text(
              'Dữ liệu giải mã:',
              style: TextStyle(fontSize: 20),
            ),
            SizedBox(height: 10),
            Text(decryptedData),
            SizedBox(height: 20),
            Text(
              'Dữ liệu giải mã:',
              style: TextStyle(fontSize: 20),
            ),
            SizedBox(height: 10),
            Text(checkData.toString()),
          ],
        ),
      ),
    );
  }

  List<String> splitData(String inputData) {
    List<String> result = inputData.split(',');
    return result.map((item) => item.trim()).toList();
  }

  Future<void> sendDecryptRequest(String hash, String aesM, String privateKey, String cipherKey) async {
    final apiUrl = 'http://localhost:8080/api/mahoadulieu/giaima';
    final requestBody = {
      'hash': hash,
      'aesM': aesM,
      'privateKey': privateKey,
      'cipherKey': cipherKey,
    };

    final response = await http.post(Uri.parse(apiUrl), body: json.encode(requestBody));
    if (response.statusCode == 200) {
      final result = json.decode(response.body);
      setState(() {
        decryptedData = result['result'] as String;
        checkData = result['checkData'] as bool;
      });
    } else {
      throw Exception('Failed to send decrypt request.');
    }
  }
}
