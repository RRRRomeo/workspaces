import numpy as np
import pandas as pd


# Generate two 3x3 random matrices
mat1 = np.random.rand(3, 3)
mat2 = np.random.rand(3, 3)

# Calculate the dot product of the two matrices using numpy.matmul()
dot_product = np.matmul(mat1, mat2)

print("Matrix1:\n", mat1)
print("\nMatrix2:\n", mat2)
print("\nDot product:\n", dot_product)

