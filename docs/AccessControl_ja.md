# Access Control

## Overview

COSMO has two access control technique, Kubernetes RBAC and COSMO DevTeam Management.

全てのCOSMOのリソースはKubernetes Custom Resource Definitionで定義されているため、Kubernertes RBACにてアクセス制御されます。
All of resources of COSMO is defined as Kubernetes Custom Resource Definition so the source of resoruce is managed by Kubernetes RBAC.

さらに、開発環境はユーザーごとにKubernertes Namespaceが分けれられ、開発者はNamespaceを自分自身の開発環境として扱うことができます。

一方で、COSMOは開発者および開発チームに対してより開発チームのマネジメントに適したアクセス制御を提供します。それがCOSMO DevTeam Managementです。

COSMO DevTeam Managementを使用して複数の開発者をチームに所属させ、またチームマネージャーにてチームメンバーを管理できます。

チームにできることは以下の通りです。
- チーム専用の開発環境テンプレートの使用
- チーム専用のユーザー設定
- チームのユーザー管理

COSMO DevTeam Managementでは、各ユーザーは以下の3つのロールに分類されます。
- `cosmo-admin`
  全チームの権限を含むシステム全体の権限を持つ管理者権限

- `manager@TEAM`
  Teamのマネジメントができるマネージャー権限
  具体的には以下の権限があります。
  - `developer@TEAM`ロールのみのユーザー作成
  - 既存ユーザーへのTeamーxロールの割り当て
  - 既存ユーザーのTeam-xロールの剥奪
  - team-xロールのみのユーザー削除
  - developer権限

- `developer@TEAM`
  チームに所属する開発者ロール
  - Team専用のテンプレートを使用することができます

## COSMO DevTeam Management

cosmoctlまたはDashboardより、`cosmo-admin`権限ロールのユーザーで`Team`を作成し、管理者を登録します。

```sh
cosmoctl login
> Login succeeded

cosmoctl whoami
> NAME          ROLE       
> dumbledore    cosmo-admin

cosmoctl create team gryffindor
> Successfuly created team "gryffindor"

cosmoctl team update-members --team "gryffindor" --users hermione-granger --manager
> Application submited to the team manager. Please wait or notify the manager to approve.

```

team

```sh
cosmoctl whoami
> NAME           ROLE       
> harry-potter   

cosmoctl join team --team gryffindor
> Application submited to the team manager. Please wait or notify the manager to approve.
```

```sh
cosmoctl whoami
> NAME                ROLE       
> hermione-granger    manager@gryffindor

cosmoctl team list-pending-approvals gryffindor
> APPLICATION-ID                       USER           ROLE
> gryffindor-developer-harry-potter    harry-potter   developer

cosmoctl team approve gryffindor-developer-harry-potter
> Successfuly approved "gryffindor-developer-harry-potter"
```

### Internal

User

```yaml
apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
kind: User
metadata:
  name: dumbledore
spec:
  role: cosmo-admin
---
apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
kind: User
metadata:
  name: hermione-granger
spec:
  role: manager@gryffindor
---
apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
kind: User
metadata:
  name: ron-weasley
spec:
  role: developer@gryffindor
```

Team

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cosmo-teams
  namespace: cosmo-system
data:
  cosmo-teams.yaml: |
    apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
    kind: TeamList
    items:
      - apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
        kind: Team
        metadata:
          name: gryffindor
        members:
          - name: harry-potter
          - name: ron-weasley
          - name: hermione-granger

      - apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
        kind: Team
        metadata:
          name: slytherin
        members:
          - name: draco-malfoy
```

Pending approval

User spec has role of developer@TEAM but it is not set in `cosmo-teams` ConfigMap

```yaml
apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
kind: User
metadata:
  name: harry-potter
spec:
  role: developer@gryffindor # set
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cosmo-teams
  namespace: cosmo-system
data:
  cosmo-teams.yaml: |
    apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
    kind: TeamList
    items:
      - apiVersion: workspace.cosmo-workspace.github.io/v1alpha1
        kind: Team
        metadata:
          name: gryffindor
        members:
          # not set
          - name: ron-weasley
          - name: hermione-granger
```

Team Template
gryffindorが所有し、ravenclawも使用可能

```yaml
apiVersion: cosmo.cosmo-workspace.github.io/v1alpha1
kind: Template
metadata:
  name: sword-of-gryffindor
  annotations:
    team.cosmo-workspace/owner-team: gryffindor
    team.cosmo-workspace/available-teams: ravenclaw
  label:
    template.cosmo-workspace/type: workspace
```

```yaml
apiVersion: cosmo.cosmo-workspace.github.io/v1alpha1
kind: ClusterTemplate
metadata:
  name: east-tower
  annotations:
    team.cosmo-workspace/owner-team: gryffindor
    useraddon.template.cosmo-workspace/default: "true"
  label:
    template.cosmo-workspace/type: user-addon
```