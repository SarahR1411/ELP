module TcTurtleParser exposing (parseTurtleProgram)

import Parser exposing (Parser, succeed, symbol, spaces, int, oneOf, (|=), (|.), loop, Step(..), lazy)
import Turtle exposing (Instruction(..))

parseTurtleProgram : Parser (List Instruction)
parseTurtleProgram =
    bracketed (lazy (\() -> loop [] repeatParser))

bracketed : Parser a -> Parser a
bracketed p =
    succeed identity
        |. symbol "["
        |. spaces
        |= p
        |. spaces
        |. symbol "]"

repeatParser : List Instruction -> Parser (Step (List Instruction) (List Instruction))
repeatParser acc =
    oneOf
        [ succeed (\instr -> Loop (acc ++ [ instr ]))
            |= instructionParser
            |. optionalComma
        , succeed (Done acc)
        ]

optionalComma : Parser ()
optionalComma =
    oneOf
        [ symbol "," |. spaces   -- Allows commas between instructions
        , spaces                 -- Also allows just spaces (no comma)
        ]


instructionParser : Parser Instruction
instructionParser =
    oneOf
        [ succeed Forward
            |. keyword "Forward"
            |= int
            |. spaces
        , succeed Left
            |. keyword "Left"
            |= int
            |. spaces
        , succeed Right
            |. keyword "Right"
            |= int
            |. spaces
        , succeed Repeat
            |. keyword "Repeat"
            |= int
            |. spaces
            |= bracketed (lazy (\() -> loop [] repeatParser))
            |. spaces
        ]

keyword : String -> Parser ()
keyword str =
    symbol str |. spaces