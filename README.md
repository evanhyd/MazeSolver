# MazeSolver
A maze solver that scans a maze image and generate solution as gif or png image.
A proof of concept that computer science can be used in real life application other than competitive programming (sarcastic).  
Good for cheating in online maze puzzle solving competition when internet is not available.  
  
# Usage  
 - Open a terminal and run:  
     `./mazesolver_(win|mac|linux).com.exe [input file] [output file] [duration] [space color] [block color] [source color] [destination color] [path color]`  
 - `duration: gif animation in seconds`  
 - `color: R,G,B from 0 - 255, separated by a comma`  
  
# Example  
 - Scan the maze image "input.png", and generate "output.gif" that's 3 seconds long.
 - Empty  pixel has a color of (255,255,255) -> White
 - Block  pixel has a color of (0,0,0)       -> Black
 - Source pixel has a color of (255,0,0)     -> Red
 - Ending pixel has a color of (0,255,0)     -> Green
 - Path   pixel has a color of (255, 0, 0)   -> Red
 - `./mazesolver_win.com.exe input.png output.gif 3 255,255,255 0,0,0 255,0,0 0,255,0 255,0,0`
 ## input.png (red spot to green spot)
 ![input.png](https://imgur.com/bLXYkNc.png)  
   
 ## output.gif
 ![output.png](https://imgur.com/w9S0100.gif)  
   
 ## More Example
 [input.png](https://i.imgur.com/yeqVJWe.png?raw=true)  
 [output.gif](https://i.imgur.com/jF2p2GU.gif?raw=true)
 
   
 # Note  
 Input  image format must be .png || .jpeg    
 Output image format must be .gif || .png  
 GIF renders at 25 frames per second.  
