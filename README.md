# Jeu de Balle Rebondissante

Un jeu simple écrit en Go.

## Description

Contrôlez une raquette pour garder la balle rebondissante en jeu. Vous avez 3 vies. Marquez des points en frappant la balle avec la raquette. Si la balle tombe en bas, vous perdez une vie. Quand toutes les vies sont épuisées, le jeu revient au menu de démarrage.

## Contrôles

- **Flèche gauche** : Déplacer la raquette à gauche
- **Flèche droite** : Déplacer la raquette à droite
- **Espace** : Démarrer le jeu (depuis le menu)

## Ressources

- `ball.png` : Une image carrée rouge de 16x16 pour la balle.
- `bounce.wav` : Un court bip joué lors du rebond.

## Comment lancer

1. Assurez-vous que Go est installé (version 1.24 ou ultérieure).
2. Clonez ou téléchargez le dépôt.
3. Ouvrez un shell dans le dossier du projet.
4. Lancez `go run main.go` ou `go build` puis `./gogame`.

La fenêtre du jeu s'ouvrira affichant le menu de démarrage. Appuyez sur Espace pour commencer à jouer.

