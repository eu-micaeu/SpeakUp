import 'dart:convert';
import 'package:http/http.dart' as http;

class AuthController {
  static const String _baseUrl = 'http://10.0.2.2:8080/user';

  static Future<Map<String, dynamic>> login(
    String email,
    String password,
  ) async {
    final url = Uri.parse('$_baseUrl/login');

    try {
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'email': email, 'password': password}),
      );

      final responseData = jsonDecode(response.body);

      if (response.statusCode == 200) {
        return responseData;
      } else {
        throw Exception(responseData['message'] ?? 'Falha no login');
      }
    } on FormatException {
      throw Exception('Resposta inválida do servidor');
    } catch (e) {
      throw Exception('Erro desconhecido: ${e.toString()}');
    }
  }
}
