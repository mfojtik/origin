package policy

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/validation/path"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
)

func buildSubjects(users, groups []string) []corev1.ObjectReference {
	subjects := []corev1.ObjectReference{}

	for _, user := range users {
		saNamespace, saName, err := serviceaccount.SplitUsername(user)
		if err == nil {
			subjects = append(subjects, corev1.ObjectReference{Kind: "ServiceAccount", Namespace: saNamespace, Name: saName})
			continue
		}

		kind := determineUserKind(user)
		subjects = append(subjects, corev1.ObjectReference{Kind: kind, Name: user})
	}

	for _, group := range groups {
		kind := determineGroupKind(group)
		subjects = append(subjects, corev1.ObjectReference{Kind: kind, Name: group})
	}

	return subjects
}

// duplicated from the user/validation package.  We need to avoid api dependencies on validation from our types.
// These validators are stable and realistically can't change.
func validateUserName(name string, _ bool) []string {
	if reasons := path.ValidatePathSegmentName(name, false); len(reasons) != 0 {
		return reasons
	}

	if strings.Contains(name, ":") {
		return []string{`may not contain ":"`}
	}
	if name == "~" {
		return []string{`may not equal "~"`}
	}
	return nil
}

// duplicated from the user/validation package.  We need to avoid api dependencies on validation from our types.
// These validators are stable and realistically can't change.
func validateGroupName(name string, _ bool) []string {
	if reasons := path.ValidatePathSegmentName(name, false); len(reasons) != 0 {
		return reasons
	}

	if strings.Contains(name, ":") {
		return []string{`may not contain ":"`}
	}
	if name == "~" {
		return []string{`may not equal "~"`}
	}
	return nil
}

func determineUserKind(user string) string {
	kind := "User"
	if len(validateUserName(user, false)) != 0 {
		kind = "SystemUser"
	}
	return kind
}

func determineGroupKind(group string) string {
	kind := "Group"
	if len(validateGroupName(group, false)) != 0 {
		kind = "SystemGroup"
	}
	return kind
}

// stringSubjectsFor returns users and groups for comparison against user.Info.  currentNamespace is used to
// to create usernames for service accounts where namespace=="".
func stringSubjectsFor(currentNamespace string, subjects []corev1.ObjectReference) ([]string, []string) {
	// these MUST be nil to indicate empty
	var users, groups []string

	for _, subject := range subjects {
		switch subject.Kind {
		case "ServiceAccount":
			namespace := currentNamespace
			if len(subject.Namespace) > 0 {
				namespace = subject.Namespace
			}
			if len(namespace) > 0 {
				users = append(users, serviceaccount.MakeUsername(namespace, subject.Name))
			}

		case "User", "SystemUser":
			users = append(users, subject.Name)

		case "Group", "SystemGroup":
			groups = append(groups, subject.Name)
		}
	}

	return users, groups
}
