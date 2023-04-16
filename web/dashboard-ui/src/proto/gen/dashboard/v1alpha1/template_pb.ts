//
//Cosmo Dashboard API
//Manipulate cosmo dashboard resource API

// @generated by protoc-gen-es v1.2.0 with parameter "target=ts"
// @generated from file dashboard/v1alpha1/template.proto (package dashboard.v1alpha1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message dashboard.v1alpha1.TemplateRequiredVars
 */
export class TemplateRequiredVars extends Message<TemplateRequiredVars> {
  /**
   * @generated from field: string var_name = 1;
   */
  varName = "";

  /**
   * @generated from field: string default_value = 2;
   */
  defaultValue = "";

  constructor(data?: PartialMessage<TemplateRequiredVars>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "dashboard.v1alpha1.TemplateRequiredVars";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "var_name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "default_value", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): TemplateRequiredVars {
    return new TemplateRequiredVars().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): TemplateRequiredVars {
    return new TemplateRequiredVars().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): TemplateRequiredVars {
    return new TemplateRequiredVars().fromJsonString(jsonString, options);
  }

  static equals(a: TemplateRequiredVars | PlainMessage<TemplateRequiredVars> | undefined, b: TemplateRequiredVars | PlainMessage<TemplateRequiredVars> | undefined): boolean {
    return proto3.util.equals(TemplateRequiredVars, a, b);
  }
}

/**
 * @generated from message dashboard.v1alpha1.Template
 */
export class Template extends Message<Template> {
  /**
   * @generated from field: string name = 1;
   */
  name = "";

  /**
   * @generated from field: string description = 2;
   */
  description = "";

  /**
   * @generated from field: repeated dashboard.v1alpha1.TemplateRequiredVars required_vars = 3;
   */
  requiredVars: TemplateRequiredVars[] = [];

  /**
   * @generated from field: optional bool is_default_user_addon = 4;
   */
  isDefaultUserAddon?: boolean;

  /**
   * @generated from field: bool is_cluster_scope = 5;
   */
  isClusterScope = false;

  constructor(data?: PartialMessage<Template>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "dashboard.v1alpha1.Template";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "description", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "required_vars", kind: "message", T: TemplateRequiredVars, repeated: true },
    { no: 4, name: "is_default_user_addon", kind: "scalar", T: 8 /* ScalarType.BOOL */, opt: true },
    { no: 5, name: "is_cluster_scope", kind: "scalar", T: 8 /* ScalarType.BOOL */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Template {
    return new Template().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Template {
    return new Template().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Template {
    return new Template().fromJsonString(jsonString, options);
  }

  static equals(a: Template | PlainMessage<Template> | undefined, b: Template | PlainMessage<Template> | undefined): boolean {
    return proto3.util.equals(Template, a, b);
  }
}

