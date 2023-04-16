//
//Cosmo Dashboard API
//Manipulate cosmo dashboard resource API

// @generated by protoc-gen-es v1.0.0 with parameter "target=ts"
// @generated from file dashboard/v1alpha1/auth_service.proto (package dashboard.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, Timestamp } from "@bufbuild/protobuf";

/**
 * @generated from message dashboard.v1alpha1.LoginRequest
 */
export class LoginRequest extends Message<LoginRequest> {
  /**
   * @generated from field: string user_name = 1;
   */
  userName = "";

  /**
   * @generated from field: string password = 2;
   */
  password = "";

  constructor(data?: PartialMessage<LoginRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "dashboard.v1alpha1.LoginRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "user_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "password", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LoginRequest {
    return new LoginRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LoginRequest {
    return new LoginRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LoginRequest {
    return new LoginRequest().fromJsonString(jsonString, options);
  }

  static equals(a: LoginRequest | PlainMessage<LoginRequest> | undefined, b: LoginRequest | PlainMessage<LoginRequest> | undefined): boolean {
    return proto3.util.equals(LoginRequest, a, b);
  }
}

/**
 * @generated from message dashboard.v1alpha1.LoginResponse
 */
export class LoginResponse extends Message<LoginResponse> {
  /**
   * @generated from field: string user_name = 1;
   */
  userName = "";

  /**
   * @generated from field: google.protobuf.Timestamp expire_at = 2;
   */
  expireAt?: Timestamp;

  /**
   * @generated from field: bool require_password_update = 3;
   */
  requirePasswordUpdate = false;

  constructor(data?: PartialMessage<LoginResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "dashboard.v1alpha1.LoginResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "user_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "expire_at", kind: "message", T: Timestamp },
    { no: 3, name: "require_password_update", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): LoginResponse {
    return new LoginResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): LoginResponse {
    return new LoginResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): LoginResponse {
    return new LoginResponse().fromJsonString(jsonString, options);
  }

  static equals(a: LoginResponse | PlainMessage<LoginResponse> | undefined, b: LoginResponse | PlainMessage<LoginResponse> | undefined): boolean {
    return proto3.util.equals(LoginResponse, a, b);
  }
}

/**
 * @generated from message dashboard.v1alpha1.VerifyResponse
 */
export class VerifyResponse extends Message<VerifyResponse> {
  /**
   * @generated from field: string user_name = 1;
   */
  userName = "";

  /**
   * @generated from field: google.protobuf.Timestamp expire_at = 2;
   */
  expireAt?: Timestamp;

  /**
   * @generated from field: bool require_password_update = 3;
   */
  requirePasswordUpdate = false;

  constructor(data?: PartialMessage<VerifyResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime = proto3;
  static readonly typeName = "dashboard.v1alpha1.VerifyResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "user_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "expire_at", kind: "message", T: Timestamp },
    { no: 3, name: "require_password_update", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): VerifyResponse {
    return new VerifyResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): VerifyResponse {
    return new VerifyResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): VerifyResponse {
    return new VerifyResponse().fromJsonString(jsonString, options);
  }

  static equals(a: VerifyResponse | PlainMessage<VerifyResponse> | undefined, b: VerifyResponse | PlainMessage<VerifyResponse> | undefined): boolean {
    return proto3.util.equals(VerifyResponse, a, b);
  }
}

