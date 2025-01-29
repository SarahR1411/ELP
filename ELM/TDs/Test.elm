module Test exposing (..) 



--première ligne module Test exposing (..). C'est là que tu définiras tes fonctions.



-- Définis la fonction addElemInList qui ajoute un élément donné, un nombre de fois donné, dans une liste donnée.
addElemInList : a -> Int -> List a -> List a -- Put "a" instead of Int so that function accept any type of elem (bool, float, char and not just int)
addElemInList elem nb list = 
    List.repeat nb elem ++ list -- creates a new list by repeating a given element nb times and then appending it to another list, list. Let's break this down step by step.



-- Définis la fonction dupli qui duplique les éléments d'une liste donnée.
dupli : List a -> List a -- Put "a" instead of Int so that function accept any type of elem (bool, float, char and not just int)
dupli list = 
    List.concatMap (\x -> [x, x]) list -- maps a function over a list and concatenates the results. The function takes an element x and returns a list containing x repeated twice.



-- Définis la fonction compress qui supprime les copies consécutives des éléments d'une liste
compress : List a -> List a -- Put "a" instead of Int so that function accept any type of elem (bool, float, char and not just int)
compress list = 
    List.foldr (\x acc -> if List.head acc == Just x then acc else x :: acc) [] list -- folds a function over a list from right to left. The function takes an element x and an accumulator acc. If the head of the accumulator is equal to x, it returns the accumulator; otherwise, it prepends x to the accumulator.

