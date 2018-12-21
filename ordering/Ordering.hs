module Ordering where

f1 :: Int -> Int
f1 n | n < 0 = n 
     | otherwise = f2 n 

f2 :: Int -> Int
f2 n = f1 (n-1)
