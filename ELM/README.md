# Turtle Graphics Elm

This project is an interactive Turtle Graphics implementation built with Elm. It allows users to input Turtle commands to draw shapes and patterns directly in the browser. The app also supports animation, speed control, and pen color/width adjustments.

## Features

- **Turtle Commands**: Draw shapes by inputting commands like `Forward`, `Left`, `Right`, `Repeat`, and more.
- **Animation**: Run the Turtle program step-by-step with animation, and adjust the animation speed.
- **Pen Control**: Toggle the pen up or down, change pen color, and adjust the pen width.
- **Example Commands**: Predefined example commands to help you get started with Turtle graphics, such as drawing a square, triangle, hexagon, and more.

## Installation

1. Ensure you have Elm installed on your system. If not, you can install it by following the instructions [here](https://elm-lang.org/docs/install).
   
2. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/SarahR1411/ELP.git
    ```

3. Navigate to the project directory:

    ```bash
    cd turtle-graphics-elm
    ```

4. Install the required dependencies:

    ```bash
    elm install
    ```

## Running the Project

1. Compile the Elm code into JavaScript:

    ```bash
    elm make src/Main.elm --output=main.js
    ```

2. Open the project in your browser:

    ```bash
    open index.html
    ```

   Alternatively, you can use a simple HTTP server:

   ```bash
   python -m SimpleHTTPServer 8000
   ```

   Then navigate to `http://localhost:8000` in your browser.

## Usage

- **Enter Turtle Commands**: Use the input field to enter Turtle commands like:

    ```elm
    [Repeat 4 [Forward 50, Left 90]]
    ```

- **Draw Instantly**: Click on "Draw Instantly" to see the result immediately after entering a command.

- **Start Animation**: Click on "Start Animation" to animate the drawing process, step by step. You can control the animation speed.

- **Pen Controls**: Toggle the pen up/down, change the pen color and width to customize the drawing style.

- **Example Commands**: Select from pre-defined example commands to quickly test the Turtle graphics functionality.

## Project Structure

- `src/Main.elm`: The main Elm application file.
- `src/TcTurtleParser.elm`: Contains the parser for interpreting Turtle commands.
- `src/Turtle.elm`: Handles the logic for executing Turtle instructions and managing the turtle's state.
- `src/View.elm`: Contains the HTML rendering functions, including the display of the canvas and UI controls.
- `index.html`: The HTML file that loads the Elm app.
