module TD1 exposing (..)

-- Pattern matching
-- Pattern matching désigne le mécanisme par lequel les types structurés peuvent être déstructurés dans trois cas :
    -- let ... in ..., qui sert à poser des définitions comme en math,
    -- case ... of ... -> ..., qui sert habituellement à distinguer des cas,
    -- la définition de fonctions.

couleur = (128, 0, 0)

let (r,v,b) = couleur in r
-- returns, 128 : number

case couleur of (r,v,b) -> r
-- returns, 128 : number

-- Use let ... in when you want to bind variables for reuse in a block of code.
-- Use case ... of when you want to handle a value directly in a single matching expression.

getRedChannel (r,v,b) = r
<function> : ( a, b, c ) -> a

getRedChannel couleur
-- 128 : number


-- Les listes sont construites soit explicitement, soit avec l'opérateur :: de sorte que ces trois expressions sont équivalentes :
[1,2,3,4]
1 :: [2,3,4]
1 :: 2 :: 3 :: 4 :: []

-- Pour déconstuire une liste, on utilise donc naturellement l'opérateur ::, mais seulement dans un case car il est nécessaire de prendre en compte le fait qu'une liste peut être vide :
case lst of
  [] -> ...
  (x :: xs) -> ...


-- EX1: cinq façons d'extraire la valeur du champs name de l'enregistrement suivant
person = { name="me", age="22" }
{ age = "22", name = "me" }
    : { age : String, name : String }

-- 1
let { name = n } = person in n
-- returns, "me" : String

-- 2
case person of { name = n } -> n
-- returns, "me" : String

-- 3
let { name = n, age = _ } = person in n
-- returns, "me" : String

-- 4
case person of { name = n, age = _ } -> n
-- returns, "me" : String

-- 5
getname { name, age } = name
getname person
-- returns, "me" : String






--FUNCTIONS
-- Functions are defined using the keyword let followed by the function name, a list of arguments, an equal sign, and the function body.
-- Les fonctions sont définies soit en nommant une lambda, soit sous forme équationnelle : inc = \x -> x + 1 (lambda à droite) est équivalent à inc x = x + 1 (équation).

-- Toutes les fonctions ne sont pas recursives: 
sign x = if x > 0 then 1 else if x == 0 then 0 else -1
estVide lst = case lst of
   [] -> True
   (x :: xs) -> False

-- Cependant, comme il n'y a pas de structure de contrôle itérative, les traitements répétés sont naturellement codés sous la forme de fonctions récursives. Par exemple :
len lst = case lst of
   [] -> 0                      -- If lst is the empty list ([]), return 0. This is the base case for recursion, where the list has no elements.
   (x :: xs) -> 1 + len xs      -- The function returns 1 + len xs, where:
                                -- 1 accounts for the current element (x).
                                -- len xs recursively calculates the length of the remaining list (xs).

    -- if lst is a non-empty list, it matches the pattern (x :: xs), where:
    -- x is the first element of the list (the "head").
    --xs is the rest of the list (the "tail").





-- PACKAGE LIST
-- Le package List est disponible par défaut et contient un grand nombre de fonctions très utiles. 
-- Très souvent les appeler nous dispensent d'écrire une fonction récursive. List.map, List.filter, List.foldr, List.foldl sont typiques des langages fonctionnels. 

-- List.map applique une fonction à chaque élément d'une liste: List.map : (a -> b) -> List a -> List b
-- List.filter filtre les éléments d'une liste: List.filter : (a -> Bool) -> List a -> List a
-- List.foldr et List.foldl sont des fonctions de pliage (ou réduction) de liste: List.foldr : (a -> b -> b) -> b -> List a -> b


-- Acceder à elm script! : 
  -- elm repl
    -- import List exposing (..)
    -- map (\x -> x + 1) [1,2,3] (eg)



