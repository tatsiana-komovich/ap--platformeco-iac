# lmru--devops--argocd-apps
Как добавить ресурсы в ArgoCD

# Структура (пример на lm-api-manager)
```
api-manager   #директория с названием кластера
 |_ stage     #окружение кластера
 |    |_ node
 |        |_ templates
 |        |     |_ instanceclasses.yaml         
 |        |     |_ node_groups.yaml      
 |        |_ Chart.yaml
 |_ prod
     |_ node
         |_ templates
         |     |_ instanceclasses.yaml         
         |     |_ node_groups.yaml      
         |_ Chart.yaml        
``` 
# Порядок добавления:
1. создать директории в соответствии со структурой
2. экспортировать данные из кластера kubectl get ng worker -oyaml
3. удалить лишние аннотации 
4. вставить экспортированные данные в файлы instanceclasses.yaml и node_groups.yaml
5. обратиться к Александру Крупину чтобы добавил директорию в argoCD
6. помержить изменения (тимлид или зам.тл)
7. убедиться что появились аннотации argoCD
