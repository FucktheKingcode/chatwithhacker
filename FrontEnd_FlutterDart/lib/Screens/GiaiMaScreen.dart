import 'dart:io';
import 'dart:typed_data';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class GiaiMaScreen extends StatefulWidget {
  @override
  _GiaiMaScreenState createState() => _GiaiMaScreenState();
}

class _GiaiMaScreenState extends State<GiaiMaScreen> {
  File? encryptedFile;
  File? privateKeyFile;
  File? decryptedFile;

  Future<void> decryptFile() async {
    if (encryptedFile == null || privateKeyFile == null) {
      showDialog(
        context: context,
        builder: (context) => AlertDialog(
          title: Text('Error'),
          content: Text(
              'Please attach both the encrypted file and private key file.'),
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
    final encryptedData = await encryptedFile!.readAsBytes();
    final privateKeyData = await privateKeyFile!.readAsBytes();

// TODO: Perform decryption here

    final decryptedData =
        Uint8List.fromList([1, 2, 3]); // Example decrypted data

    final decryptedFilePath = await _saveFile(decryptedData, 'Decrypted File');

    setState(() {
      decryptedFile = File(decryptedFilePath);
    });

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Decryption Successful'),
        content: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('File Data: $decryptedFilePath'),
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

  Future<void> _pickEncryptedFile() async {
    final result = await FilePicker.platform.pickFiles();
    if (result != null && result.files.isNotEmpty) {
      setState(() {
        encryptedFile = File(result.files.first.path!);
      });
    }
  }

  Future<void> _pickPrivateKeyFile() async {
    final result = await FilePicker.platform.pickFiles();
    if (result != null && result.files.isNotEmpty) {
      setState(() {
        privateKeyFile = File(result.files.first.path!);
      });
    }
  }

  Future<String> _saveFile(Uint8List data, String fileName) async {
    final tempDir = Directory.systemTemp;
    final file = File('${tempDir.path}/$fileName');
    await file.writeAsBytes(data);

    return file.path;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Decrypt Screen'),
      ),
      body: Column(
        children: [
          ListTile(
            title: Text('File Mã Hóa'),
            subtitle: encryptedFile == null
                ? Text('No file selected')
                : Text(encryptedFile!.path),
            trailing: IconButton(
              icon: Icon(Icons.attach_file),
              onPressed: _pickEncryptedFile,
            ),
          ),
          ListTile(
            title: Text('File Private Key'),
            subtitle: privateKeyFile == null
                ? Text('No file selected')
                : Text(privateKeyFile!.path),
            trailing: IconButton(
              icon: Icon(Icons.attach_file),
              onPressed: _pickPrivateKeyFile,
            ),
          ),
          ElevatedButton(
            child: Text('Giải Mã'),
            onPressed: decryptFile,
          ),
          if (decryptedFile != null)
            ListTile(
              title: Text('File Dữ liệu'),
              subtitle: Text(decryptedFile!.path),
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
