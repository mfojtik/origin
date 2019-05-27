package api

import (
	"regexp"

	"k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ValidateProvisionRequest(preq *ProvisionRequest) field.ErrorList {
	errors := ValidateUUID(field.NewPath("service_id"), ServiceID)
	errors = append(errors, ValidateUUID(field.NewPath("plan_id"), PlanID)...)

	if Platform == "" {
		errors = append(errors, field.Required(field.NewPath("context.platform"), ""))
	} else if Platform != ContextPlatformKubernetes {
		errors = append(errors, field.Invalid(field.NewPath("context.platform"), Platform, "must equal "+ContextPlatformKubernetes))
	}

	if Namespace == "" {
		errors = append(errors, field.Required(field.NewPath("context.namespace"), ""))
	} else {
		for _, msg := range validation.ValidateNamespaceName(Namespace, false) {
			errors = append(errors, field.Invalid(field.NewPath("context.namespace"), Namespace, msg))
		}
	}

	return errors
}

func ValidateBindRequest(breq *BindRequest) field.ErrorList {
	errors := ValidateUUID(field.NewPath("service_id"), ServiceID)
	errors = append(errors, ValidateUUID(field.NewPath("plan_id"), PlanID)...)

	return errors
}

var uuidRegex = regexp.MustCompile("^(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")

func ValidateUUID(path *field.Path, uuid string) field.ErrorList {
	if uuidRegex.MatchString(uuid) {
		return nil
	}
	return field.ErrorList{field.Invalid(path, uuid, "must be a valid UUID")}
}
