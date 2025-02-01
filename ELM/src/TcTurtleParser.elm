module TcTurtleParser exposing (parseTurtleProgram)

-- Import the Parser module, which provides functions for parsing strings into values
import Parser exposing (Parser, succeed, symbol, spaces, int, oneOf, (|=), (|.), loop, Step(..), lazy)

-- Import the Instruction type from the Turtle module, which represents turtle commands
import Turtle exposing (Instruction(..))

-- The parseTurtleProgram function parses a turtle program from a string
parseTurtleProgram : Parser (List Instruction)
parseTurtleProgram =
    -- Parse the program as a bracketed expression, which starts with "[" and ends with "]"
    -- The program is recursively parsed using the repeatParser function
    bracketed (lazy (\() -> loop [] repeatParser))

-- Helper function to parse bracketed expressions (i.e., anything between "[" and "]")
bracketed : Parser a -> Parser a
bracketed p =
    -- Parse a "[" symbol, followed by spaces, then the inner parser, and finally a "]" symbol
    succeed identity
        |. symbol "["   -- Match the opening bracket
        |. spaces       -- Allow spaces after the opening bracket
        |= p            -- Parse the inner expression
        |. spaces       -- Allow spaces before the closing bracket
        |. symbol "]"   -- Match the closing bracket

-- The repeatParser function is responsible for parsing the instructions inside a Repeat block
repeatParser : List Instruction -> Parser (Step (List Instruction) (List Instruction))
repeatParser acc =
    oneOf
        [ -- Parse a regular instruction and add it to the list of instructions
          -- If we parse an instruction, it is added to the list and wrapped in a Loop step
          succeed (\instr -> Loop (acc ++ [ instr ]))
            |= instructionParser    -- Use instructionParser to parse individual instructions
            |. optionalComma        -- Allow optional commas between instructions
        ,                           -- If we've reached the end of the Repeat block, return the accumulated instructions
          succeed (Done acc)
        ]

-- Parser to allow optional commas between instructions (or just spaces)
optionalComma : Parser ()
optionalComma =
    oneOf
        [ symbol "," |. spaces   -- Match a comma followed by spaces
        , spaces                 -- Or just match spaces with no comma
        ]


-- The instructionParser function parses individual turtle instructions like Forward, Left, Right, and Repeat
instructionParser : Parser Instruction
instructionParser =
    oneOf
        [ -- Parse the "Forward" instruction followed by an integer value (distance)
          succeed Forward
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

-- The keyword function matches a specific string followed by spaces (used for matching "Forward", "Left", etc.)
keyword : String -> Parser ()
keyword str =
    symbol str |. spaces    -- Match the given string and allow spaces after it