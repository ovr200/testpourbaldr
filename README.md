Api rest pour le test
===================


----------
> **Fait:**

> - Utilisateur stocké en BDD (mysql / pass crypté , décrypté à sa récéption).
> - Connexion avec JWT (la route /login sert le token après vérification en BDD)
> - **/signup** créer un compte et envoie directement un email contenant l'email est le mot de passe.
> - **/reset** renvoie un nouveau mot de passe par email à l'authentifié.
> - **/album/:pagination** liste albums (pas d'authentification nécessaire).
> - **/favorite/:pagination** liste les favoris de l'utilisateur authentifié.
>  - **/favorite/:pagination/:iduser** liste les favoris de l'utilisateur ciblé (:iduser).
> - **/favorite/:IdAlbum** ajoute/retire des favoris l'album.

-------------
> **<i class="icon-pencil"></i>A faire:**

> - Validation des données reçues
> -  Faire des fonctions pour éviter les répétitions pour la partie JWT (check , claims)
> -  Actions / Erreurs
> -  Cache avec redis ( )
> - Tous recommencer proprement

#### <i class="icon-file"></i> Postman

J'ai joint le fichier postman que j'ai utilisé pour mes tests.