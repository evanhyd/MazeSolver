# MazeSolver
A maze solver that scans a maze image and generate solution as an animated gif file.  
A proof of concept that computer science can be used in real life application other than competitive programming (sarcastic).  
Good for cheating in online maze puzzle solving competition.  

# Usage  
 - Open a terminal and run:  
     `./mazesolver.com.exe [input file] [output file] [duration] [space color] [block color] [source color] [destination color] [path color]`  
 - `duration: gif animation in seconds`  
 - `color: R,G,B from 0 - 255, separated by a comma`  

# Example  
 - Scan the maze image "input.png", and generate "output.gif" that's 1 seconds long.
 - Empty  pixel has a color of (255,255,255) -> White
 - Block  pixel has a color of (0,0,0)       -> Black
 - Source pixel has a color of (255,0,0)     -> Red
 - Ending pixel has a color of (0,255,0)     -> Green
 - Path   pixel has a color of (255, 0, 0)   -> Red
 - `./mazesolver.com.exe input.png output.gif 5 255,255,255 0,0,0 255,0,0 0,255,0 255,0,0`
 ## input.png (top left red pixel to bottom right green pixel)
 ![input.png](https://imgur.com/ujsHviG.png)  
 
 ## output.gif
 ![output.png](https://imgur.com/6QZARvT.gif)  
 
 ## Even bigger?
 [input.png](https://i.imgur.com/yeqVJWe.png?raw=true)  
 [output.gif](https://i.imgur.com/jF2p2GU.gif?raw=true)
 
 
 # Side Note  
 Input  image format must be .png  
 Output image format must be .gif || .png  
 GIF renders at 50 frames per second.
