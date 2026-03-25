package mdht.rego

default auth := false

auth if {
	jwt_valid
	jwt_payload.iss == input.ISS
}

jwt_valid := valid if {
	[valid, header, payload] := verify_jwt
}

jwt_payload := payload if {
	[valid, header, payload] := verify_jwt
	valid
}

verify_jwt := io.jwt.decode_verify(input.Token, {"cert": input.Key})
