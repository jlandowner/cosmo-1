// @generated by protoc-gen-connect-web v0.3.1 with parameter "target=ts"
// @generated from file auth-proxy/v1alpha1/authproxy.proto (package authproxy.v1alpha1, syntax proto3)
/* eslint-disable */
/* @ts-nocheck */

import {Empty, LoginRequest} from "./authproxy_pb.js";
import {MethodKind} from "@bufbuild/protobuf";

/**
 * @generated from service authproxy.v1alpha1.AuthProxyService
 */
export const AuthProxyService = {
  typeName: "authproxy.v1alpha1.AuthProxyService",
  methods: {
    /**
     * @generated from rpc authproxy.v1alpha1.AuthProxyService.Login
     */
    login: {
      name: "Login",
      I: LoginRequest,
      O: Empty,
      kind: MethodKind.Unary,
    },
  }
} as const;

