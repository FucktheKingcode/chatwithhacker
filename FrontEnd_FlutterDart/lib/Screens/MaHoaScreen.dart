import 'dart:io';
import 'dart:typed_data';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class MaHoaScreen extends StatefulWidget {
  @override
  _MaHoaScreenState createState() => _MaHoaScreenState();
}

class _MaHoaScreenState extends State<MaHoaScreen> {
  File? attachmentFile1;
  File? attachmentFile2;

  Future<void> encryptFiles() async {
    if (attachmentFile1 == null || attachmentFile2 == null) {
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Text('Error'),
          content: Text('Please attach both files.'),
          actions: [
            TextButton(
              child: Text('OK'),
              onPressed: () => Navigator.pop(context),
            ),
          ],
        ),
      );
      return;
    }

    final file1Bytes = await attachmentFile1!.readAsBytes();
    final file2Bytes = await attachmentFile2!.readAsBytes();

    // TODO: Perform encryption and hash calculations here

    final encryptedData =
        Uint8List.fromList([1, 2, 3]); // Example encrypted data
    final hashData = Uint8List.fromList([4, 5, 6]); // Example hash data

    // Save encrypted file
    final encryptedFilePath = await _saveFile(encryptedData, 'encrypted_file');

    // Save hash file
    final hashFilePath = await _saveFile(hashData, 'hash_file');

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Mã hóa thành công'),
        content: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('File mã hóa : $encryptedFilePath'),
            Text('File chứa hàm băm của file data: $hashFilePath'),
          ],
        ),
        actions: [
          TextButton(
            child: Text('OK'),
            onPressed: () => Navigator.pop(context),
          ),
        ],
      ),
    );
  }

  Future<String> _saveFile(Uint8List data, String fileName) async {
    final tempDir = Directory.systemTemp;
    final file = File('${tempDir.path}/$fileName');

    await file.writeAsBytes(data);

    return file.path;
  }

  Future<void> _pickFile(int fileNumber) async {
    final result = await FilePicker.platform.pickFiles();
    if (result != null && result.files.isNotEmpty) {
      setState(() {
        if (fileNumber == 1) {
          attachmentFile1 = File(result.files.first.path!);
        } else if (fileNumber == 2) {
          attachmentFile2 = File(result.files.first.path!);
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Ma Hoa Screen'),
      ),
      body: Column(
        children: [
          ListTile(
            title: Text('File Dữ Liệu'),
            subtitle: attachmentFile1 == null
                ? Text('No file selected')
                : Text(attachmentFile1!.path),
            trailing: IconButton(
              icon: Icon(Icons.attach_file),
              onPressed: () => _pickFile(1),
            ),
          ),
          ListTile(
            title: Text('File Public Key'),
            subtitle: attachmentFile2 == null
                ? Text('No file selected')
                : Text(attachmentFile2!.path),
            trailing: IconButton(
              icon: Icon(Icons.attach_file),
              onPressed: () => _pickFile(2),
            ),
          ),
          ElevatedButton(
            child: Text('Mã Hóa'),
            onPressed: encryptFiles,
          ),
          if (attachmentFile1 != null)
            ListTile(
              title: Text('File Mã Hóa'),
              subtitle: Text(attachmentFile1!.path),
              trailing: IconButton(
                icon: Icon(Icons.file_download),
                onPressed: () {
// TODO: Implement file download functionality
                },
              ),
            ),
          if (attachmentFile1 != null)
            ListTile(
              title: Text('File Mã Băm'),
              subtitle: Text(attachmentFile1!.path),
              trailing: IconButton(
                icon: Icon(Icons.file_download),
                onPressed: () {
// TODO: Implement file download functionality
                },
              ),
            ),
        ],
      ),
    );
  }
}
