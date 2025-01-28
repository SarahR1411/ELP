module TD2 exposing (..)
import Browser.Dom exposing (Error)

-- Les types algébriques, appelés custom types dans le guide, sont une combinaison de :

-- type somme (ou encore union, enum, etc.), comme type Couleur = Rouge | Noir,
-- type produit (qui ne sont rien d'autres que des tuples personnalisés), comme type Point = DonnePoint Float Float,et récursivité.



--Avant de donner plus d'exemples, mettons-nous d'accord sur les mots. A gauche du =, c'est le monde des types. 
--On a un constructeur de type. A droite du =, on a le monde des valeurs. On a une liste de constructeurs de valeurs séparés par des |. 
--Un constructeur de valeur peut avoir zéro paramètre (Rouge, Noir) ou plusieurs paramètres (Point Float Float a deux paramètres, chacun de type Float). 
--Pour obtenir une valeur, on appelle le constructeur de valeur :
--Tous les constructeurs commencent par une majuscule.

type Couleur = Rouge | Noir
-- Couleur : Type
type Point = Point Float Float
-- > Point 3.2 5.4
-- Point 3.2 5.4 : Point

-- other eg
type Maybe a = Just a | Nothing
-- > Just 'c'
-- Just 'c' : Maybe Char
-- > Just 5.2
-- Just 5.2 : Maybe Float
-- > Nothing
-- Nothing : Maybe a

type Result error value = Ok value | Err error
-- > Ok "it works!"
-- Ok ("it works!") : Result error String
-- > Err "it fails!"
-- Err ("it fails!") : Result String value










-- Ex1
-- Donne, pour chacun des cinq types précédents (Couleur, Point, Maybe, Result, StackInt), le constructeur de type, ses paramètres, la liste des constructeurs de valeurs et leurs paramètres.
type Couleur = Rouge | Noir
-- type Couleur : Type
-- constructeur de type : Couleur
-- paramètres : aucun
-- liste des constructeurs de valeurs : Rouge, Noir

type Point = Point Float Float
-- type Point: type
-- constructeur de type : Point 
-- paramètres: Float, Float
-- liste des constructeurs de valeurs: Point Float Float

type Maybe a = Just a | Nothing
-- type Maybe a: type
-- constructeur de type: Maybe
-- paramètres: a
-- liste des constructeurs de valeurs: Just a, Nothing

type Result error value = OK value | Err error
-- type Result error value: type
-- constructeur de type: Result
-- paramètres: error, value
-- liste des constructeurs de valeurs: OK value, Err error

type StackInt = Empty | Push Int StackInt
-- type StackInt: type
-- constructor de type: StackInt
-- paramètres: Int
-- liste des constructeurs de valeurs: Empty, Push Int StackInt





--EX2:

-- Propose un type CouleurCarte pour modéliser la couleur (ou enseigne) d'une carte à jouer.
type CouleurCarte = Coeur | Carreau | Pique | Treffle
-- type CouleurCarte : Type
-- constructeur de type : CouleurCarte
-- paramètres : aucun
-- liste des constructeurs de valeurs : Coeur, Carreau, Pique, Treffle

-- Propose un type ValeurCarte pour modéliser la valeur d'une carte à jouer.
type ValeurCarte = Deux | Trois | Quatre | Cinq | Six | Sept | Huit | Neuf | Dix | Valet | Dame | Roi | As
-- type ValeurCarte : Type
-- constructeur de type : ValeurCarte
-- paramètres : aucun
-- liste des constructeurs de valeurs : Deux, Trois, Quatre, Cinq, Six, Sept, Huit, Neuf, Dix, Valet, Dame, Roi, As

-- Propose un type Carte pour modéliser une carte à jouer.
type Carte = Carte CouleurCarte ValeurCarte
-- type Carte : Type
-- constructeur de type : Carte
-- paramètres : CouleurCarte, ValeurCarte
-- liste des constructeurs de valeurs : Carte CouleurCarte ValeurCarte

-- Crée la carte correspondant à l'as de Trèfle, puis crée la liste des quatre as.
carteAsTrefle = Carte Treffle As
-- carteAsTrefle : Carte
ListeQuatreAs = [Carte Coeur As, Carte Carreau As, Carte Pique As, Carte Treffle As]
-- ListeQuatreAs : List Carte As dans les diff couleurs 

--EX3 
-- Propose un type paramétré Tree pour modéliser un arbre binaire. Puis, crée un arbre vide, ainsi qu'un arbre contenant au moins trois nombres flottants. Enfin, écris une fonction qui retourne la hauteur d'une arbre.
type Tree = Vide | Noeud Float Tree Tree    -- A proper binary tree should allow each node to have left and right subtrees. Therefore, the Tree type should look like this:
-- type Tree : Type
-- constructeur de type : Tree
-- paramètres : Float, Tree, Tree
-- liste des constructeurs de valeurs : Vide, Noeud Float Tree Tree

hauteur : Tree -> Int
hauteur tree = case tree of
    Vide -> 0
    Noeud _ g d -> 1 + max (hauteur g) (hauteur d)
-- hauteur : Tree -> number
-- hauteur Tree : 0
-- hauteur (Noeud 1 (Noeud 2 Vide Vide) (Noeud 3 Vide Vide)) : 2







-- ARCHITECTURE IN ELM
