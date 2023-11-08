import 'dart:convert';
import 'package:http/http.dart' as http;

class EncryptionResult {
  final String hash;
  final String encryptedData;
  final List<String> publicKey;
  final List<String> privateKey;
  final String encryptionKey;

  EncryptionResult({
    required this.hash,
    required this.encryptedData,
    required this.publicKey,
    required this.privateKey,
    required this.encryptionKey,
  });

  factory EncryptionResult.fromJson(Map<String, dynamic> json) {
    return EncryptionResult(
      hash: json['hash'],
      encryptedData: json['encrypted_data'],
      publicKey: List<String>.from(json['public_key']),
      privateKey: List<String>.from(json['private_key']),
      encryptionKey: json['encryption_key'],
    );
  }
}

Future<EncryptionResult> getEncryptionResult(String plaintext) async {
  final url = 'http://localhost:8080/api/mahoadulieu/$plaintext';
  final response = await http.get(Uri.parse(url));

  if (response.statusCode == 200) {
    final jsonData = json.decode(response.body);
    return EncryptionResult.fromJson(jsonData);
  } else {
    throw Exception('Failed to fetch encryption result');
  }
}
