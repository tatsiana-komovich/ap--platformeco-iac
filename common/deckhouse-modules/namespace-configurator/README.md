# Namespace configurator
## Table of Contents
1. [Description](#description)
2. [Configuration file example](#example)
3. [Custom labels](#custom_labels)
4. [F.A.Q](#faq)

## Description <a name="description"></a>
**What**: Конфигуратор неймспейсов вешает пользовательские лейблы на неймспейсы, соответствующие регулярному выражению.  

**Why**: Мы используем конфигуратор неймспейсов, чтобы привязать namespace в Kubernetes к `ProductID`. Благодаря этому, мы можем точнее распределять тех. долги и счета за ресурсы между продуктами в кластере. 

**How**: Конфигурация из ArgoCD применяется для каждого кластера индивидуально, в зависимости от названия файла в директории `values/`. Отдельно в каждом кластере конфигуарцию считывает Kubernetes-оператор. Оператор проверяет новые и существующие неймспейсы на совпадение с одним из паттернов в конфигурации и добавляет к нему соответствующий лейбл.

## Configuration file example <a name="example"></a>
Важно, чтобы регулярное выражение имело якорь начала `^` и якорь конца `$` (если мы не используем в конце паттерна`.*`), чтобы избежать потенциальных проблем с сопоставлением.  
Сам лейбл указывается "как есть" в самом начале для каждого продукта. Мы решили использовать лейбл `product-id` для продукта.
Ниже приведён пример файла для 2 продуктов с разными паттернами неймспейсов:
```yaml
---
products:
  - product-id: <product-id>
    namespace_pattern:
      - ^product-.*          # matches `product-test`, `product-prod`, etc.
    exclusions:
      - ^product-another-product$  # excludes the namespace for the product

  - product-id: <another_product-id>
    namespace_pattern:
      - ^another-product-.*  # matches `another-product-test`, `another-product-prod`, etc.
      - ^another-product$    # matches `another-product` 
  ...
```

## Custom labels <a name="custom_labels"></a>
Помимо установки лейблов `product-id` на неймспейсы есть возможность вешать и другие лейблы.
Ниже приведён пример с лейблами для Istio, Prometheus и Loki:
```yaml
istio:
  - istio-injection: enabled
    namespace_pattern:
      - ^my-namespace-.*
    exclusions:
      - ^my-namespace-but-prod$

prometheus:
  - prometheus.deckhouse.io/monitor-watcher-enabled: "true"
    namespace_pattern:
      - ^my-namespace$

loki:
  - log-shipping: enabled
    namespace_pattern:
      - ^my-namespace$
```

Ключ верхнего уровня (`istio`, `prometheus` и `loki` в примере) служит лишь для описания назначения блока лейблов и может быть любым.  
Лейблов в блоке может быть несколько. В этом случае каждый из них должен располагаться на отдельной строке.

## F.A.Q <a name="faq"></a>
### Как узнать правильное название файла для моего кластера?
У нас есть [список кластеров для ArgoCD](https://github.lmru.tech/adeo/lmru--devops--argocd/blob/master/argocd-crd/charts/clusters/values.yaml), где поле `name` является именем кластера и, соответственно, необходимым именем файла для кофигуратора неймспейсов. Например, для кластера `x1-fin-prod` имя файла должно быть `x1-fin-prod-cluster.yaml`.  

### Можно ли закрепить за одним неймспейсом несколько ProductID?
Сущности Kubernetes `label` и `annotation` не могут иметь несколько одинаковых полей.  
Поэтому в рамках конфигуратора за **одним неймспейсом может быть закреплён только один продукт**.

### Регулярные выражения для 2 разных продуктов накладываются между собой, как исправить?
Добавьте в исключения продукта регулярное выражение, которое нужно исключить. 
```yaml
---
products:
  - product-id: <product-id>
    namespace_pattern:
      - ^product-.*                # matches `product-test`, `product-prod`, etc.
    exclusions:
      - ^product-lm-.*             # excludes the namespaces for the product
```

### Что делать если у меня в одном неймспейсе живут несколько продуктов?
Развозить их по отдельным неймспейсам по возможности.  
Такую парадигму очень трудно поддерживать и накладывает некоторые ограничения:
- усложняется подсчёт потребления ресурсов за каждым продуктом;
- усложняется конфигурация в ArgoCD;
- приходится использовать исключения для правильно расчета тех. долга от ИБ;
- невозможно выставить ограничение по ресурсам на неймспейс индивидуально для продукта;
