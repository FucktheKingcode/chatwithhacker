import 'dart:io';
import 'package:file_picker/file_picker.dart';
import 'package:firebase_storage/firebase_storage.dart';
import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:uuid/uuid.dart';

class ChatRoom extends StatefulWidget {
  final Map<String, dynamic> userMap;
  final String chatRoomId;

  ChatRoom({required this.chatRoomId, required this.userMap});

  @override
  State<ChatRoom> createState() => _ChatRoomState();
}

class _ChatRoomState extends State<ChatRoom> {
  final TextEditingController _message = TextEditingController();

  final FirebaseFirestore _firestore = FirebaseFirestore.instance;

  final FirebaseAuth _auth = FirebaseAuth.instance;

  File? file;
  UploadTask? uploadTask;

  Future<void> getFile() async {
    FilePickerResult? result = await FilePicker.platform.pickFiles();

    if (result != null) {
      PlatformFile platformFile = result.files.first;
      file = File(platformFile.path!);
      uploadFile(file!);
    }
  }

  Future<void> downloadFile(String fileUrl) async {
    try {
      final Reference ref = FirebaseStorage.instance.refFromURL(fileUrl);
      final String fileName = ref.name;

      String? savePath = await FilePicker.platform.getDirectoryPath();

      if (savePath != null) {
        final File file = File('$savePath/$fileName');

        if (await file.exists()) {
          // Nếu file đã tồn tại, không cần tải lại
          print('File already exists.');
          return;
        }

        showDialog(
          context: context,
          barrierDismissible: false,
          builder: (BuildContext context) {
            return AlertDialog(
              title: Text('Downloading File'),
              content: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  CircularProgressIndicator(),
                  SizedBox(height: 16),
                  Text('Please wait...'),
                ],
              ),
            );
          },
        );

        await ref.writeToFile(file);

        Navigator.pop(context); // Đóng hộp thoại tiến trình tải

        print('File downloaded successfully.');
      } else {
        print('No directory selected.');
      }
    } catch (e) {
      print('Failed to download file: $e');
    }
  }

  Future<void> uploadFile(File file) async {
    String fileName = const Uuid().v1();
    int status = 1;

    try {
      await FirebaseFirestore.instance
          .collection('chatroom')
          .doc(widget.chatRoomId)
          .collection('chats')
          .doc(fileName)
          .set({
        "sendby": FirebaseAuth.instance.currentUser!.displayName,
        "message": "",
        "type": "file",
        "time": FieldValue.serverTimestamp(),
        "filename": fileName
      });

      final ref = FirebaseStorage.instance.ref().child('files').child(fileName);
      final metadata = SettableMetadata(contentType: 'application/pdf');
      final uploadTask = ref.putFile(file, metadata);

      await uploadTask.catchError((error) async {
        await FirebaseFirestore.instance
            .collection('chatroom')
            .doc(widget.chatRoomId)
            .collection('chats')
            .doc(fileName)
            .delete();

        status = 0;
      });

      if (status == 1) {
        TaskSnapshot snapshot = await uploadTask.whenComplete(() {});
        String fileUrl = await snapshot.ref.getDownloadURL();

        await FirebaseFirestore.instance
            .collection('chatroom')
            .doc(widget.chatRoomId)
            .collection('chats')
            .doc(fileName)
            .update({"message": fileUrl});

        print(fileUrl);
      }
    } catch (e) {
      print('Failed to upload file: $e');
    }
  }

  void onSendMessage() async {
    if (_message.text.isNotEmpty) {
      Map<String, dynamic> messages = {
        "sendby": _auth.currentUser!.displayName,
        "message": _message.text,
        "type": "text",
        "time": FieldValue.serverTimestamp(),
      };

      _message.clear();
      await _firestore
          .collection('chatroom')
          .doc(widget.chatRoomId)
          .collection('chats')
          .add(messages);
    } else {
      print("Enter Some Text");
    }
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;

    return Scaffold(
      appBar: AppBar(
        title: StreamBuilder<DocumentSnapshot>(
          stream: _firestore
              .collection("users")
              .doc(widget.userMap['uid'])
              .snapshots(),
          builder: (context, snapshot) {
            if (snapshot.hasData) {
              return Container(
                child: Column(
                  children: [
                    Text(widget.userMap['name']),
                    Text(
                      snapshot.data!['status'],
                      style: TextStyle(fontSize: 14),
                    ),
                  ],
                ),
              );
            } else {
              return Container();
            }
          },
        ),
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            Container(
              height: size.height / 1.25,
              width: size.width,
              child: StreamBuilder<QuerySnapshot>(
                stream: _firestore
                    .collection('chatroom')
                    .doc(widget.chatRoomId)
                    .collection('chats')
                    .orderBy("time", descending: false)
                    .snapshots(),
                builder: (BuildContext context,
                    AsyncSnapshot<QuerySnapshot> snapshot) {
                  if (snapshot.hasData) {
                    return ListView.builder(
                      itemCount: snapshot.data!.docs.length,
                      itemBuilder: (context, index) {
                        Map<String, dynamic> map = snapshot.data!.docs[index]
                            .data() as Map<String, dynamic>;
                        return messages(size, map, context);
                      },
                    );
                  } else {
                    return Container();
                  }
                },
              ),
            ),
            Container(
              height: size.height / 10,
              width: size.width,
              alignment: Alignment.center,
              child: Container(
                height: size.height / 12,
                width: size.width / 1.1,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Container(
                      height: size.height / 17,
                      width: size.width / 1.3,
                      child: TextField(
                        controller: _message,
                        decoration: InputDecoration(
                            suffixIcon: IconButton(
                              onPressed: () => getFile(),
                              icon: Icon(Icons.file_copy_outlined),
                            ),
                            hintText: "Send Message",
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(8),
                            )),
                      ),
                    ),
                    IconButton(
                        icon: Icon(Icons.send), onPressed: onSendMessage),
                  ],
                ),
              ),
            ),
            if (uploadTask != null)
              StreamBuilder<TaskSnapshot>(
                stream: uploadTask!.snapshotEvents,
                builder: (context, snapshot) {
                  if (snapshot.hasData) {
                    final TaskSnapshot? taskSnapshot = snapshot.data;
                    final int progress = (100 *
                            taskSnapshot!.bytesTransferred /
                            taskSnapshot.totalBytes)
                        .round();
                    return Container(
                      padding: EdgeInsets.all(8),
                      child: Text('Uploading: $progress%'),
                    );
                  } else {
                    return Container();
                  }
                },
              ),
          ],
        ),
      ),
    );
  }

  Widget messages(Size size, Map<String, dynamic> map, BuildContext context) {
    return map['type'] == "text"
        ? Container(
            width: size.width,
            alignment: map['sendby'] == _auth.currentUser!.displayName
                ? Alignment.centerRight
                : Alignment.centerLeft,
            child: Container(
              padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 14),
              margin: const EdgeInsets.symmetric(vertical: 5, horizontal: 8),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(15),
                color: Colors.blue,
              ),
              child: Text(
                map['message'],
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w500,
                  color: Colors.white,
                ),
              ),
            ),
          )
        : Container(
            width: size.width,
            alignment: map['sendby'] == _auth.currentUser!.displayName
                ? Alignment.centerRight
                : Alignment.centerLeft,
            child: Container(
              width: size.width,
              padding: EdgeInsets.symmetric(vertical: 10, horizontal: 14),
              margin: EdgeInsets.symmetric(vertical: 5, horizontal: 8),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(15),
                color: Colors.blue,
              ),
              child: InkWell(
                onTap: () => downloadFile(map['message']),
                child: Row(
                  children: [
                    Icon(Icons.insert_drive_file),
                    SizedBox(width: 8),
                    if (map['filename'] != null)
                      Text(
                        map['filename'],
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w500,
                          color: Colors.white,
                        ),
                      ),
                  ],
                ),
              ),
            ),
          );
  }
}
