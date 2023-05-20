import math
# Your data
X = [66, 68, 61, 62, 71, 71, 66, 78, 72, 64, 73, 78, 76, 72, 75, 69, 74, 79, 74, 68, 81, 61, 72, 64, 69, 80, 81, 82, 74, 71]


Y = [70, 66, 66, 60, 72, 76, 73, 82, 66, 74, 81, 82, 78, 74, 76, 74, 70, 78, 75, 64, 87, 62, 73, 56, 72, 79, 82, 86, 70, 68]

if (len(X) != len(Y)):
    print("X and Y must have the equal length.")
    exit(-1)

# Print of the correlation field
fig = px.scatter(x=X, y=Y)
fig.show()

# Needed fuctions
def sum(A):
    a = 0
    for x in A:
        a += x
    return a

def sum2(A):
    a = 0
    for x in A:
        a += (x)**2
    return a

def sumN(A: list, N: int):
    sum = 0
    for x in A:
        sum += x**N
    return sum

def M(A):
    s = sum(A)
    return s/len(A); 

def D(A):
    return sum2(A)/len(A) - M(A)**2

def sigma(A):
    return math.sqrt(D(A))

def cov(A, B):
    s = 0
    for i in range(len(A)):
        s += (A[i])*(B[i])
    return s/len(A) - M(A)*M(B)

def r_v(A, B):
    return cov(A, B)/(math.sqrt(D(A)*(D(B))))

def get_rangs(G: list):
    G.sort()
    value_and_rang = {}
    rang = 1
    count = 1
    ix = 1
    value = G[0]
    while (ix < len(G)):
        if (G[ix-1] == G[ix]):
            rang += ix+1
            count += 1
            ix += 1
        else:
            rang /= count
            value_and_rang[value] = rang
            # print(f"rang({value}) = {rang} | {count}")
            rang = ix + 1
            count = 1
            value = G[ix]
            ix += 1
    rang /= count
    value_and_rang[value] = rang
    # print(f"rang({value}) = {rang} | {count}")
    # print("---")
    return value_and_rang

def rangs_for_Kendal(X: list) -> list:
    dict_rangs = get_rangs(X.copy())
    dict_of_count_eq = {}
    l_result = []
    for x in X:
        if x in dict_of_count_eq.keys():
            dict_of_count_eq[x] += 1
        else:
            dict_of_count_eq[x] = 1
        l_result.append(dict_rangs[x] + ((-1)**(dict_of_count_eq[x]))*(dict_of_count_eq[x]-1)/100)
    return l_result

def sumXY(X, Y: list):
    s = 0
    for i in range(len(X)):
        s += X[i] * Y[i]
    return s

def sumXNY(X, Y: list, N: int):
    s = 0
    for i in range(len(X)):
        s += X[i]**N * Y[i]
    return s

def k(X, Y):
    return cov(X, Y)/D(X)

def b(X, Y):
    return M(Y) - k(X, Y)*M(X)

# Solving 
print("Решение пункта 2:")
print(f"|X| = |Y| = {len(X)}")
print(f"M(X) = {M(X)}, D(X) = {D(X)}, sigma(X) = {sigma(X)}")
print(f"M(Y) = {M(Y)}, D(Y) = {D(Y)}, sigma(Y) = {sigma(Y)}")
print("------------------------------------------------")
print(f"cov(X, Y) = {cov(X, Y)}")
print(f"r_v(X, Y) = {r_v(X, Y)}")
print("Проверка осуществляется по Стьюденту со степенью свободы df = n - 2 и уровнем значимости 0.05")
print(f"t_набл = {r_v(X, Y)*math.sqrt(len(X) -2)/math.sqrt(1 - r_v(X, Y)**2)}")
print(f"Определение значимости и нахождение доверительного интервала оставляю на тебя, мой друг")
print()

print("Решение пункта 3:\n---")
print("Найдем коэф. корреляции Спирмана:")
print("Для этого отсортируем массивы и проранжируем их:")
dict_X_values_and_rangs = get_rangs(X.copy())
dict_Y_values_and_rangs = get_rangs(Y.copy())
print("Вывод рангов:\n---")
print("Для X:")
for x in sorted(X):
    print(f"rang({x}) = {dict_X_values_and_rangs[x]}")
print("---")
print("Для Y:")
for y in sorted(Y):
    print(f"rang({y}) = {dict_Y_values_and_rangs[y]}")
arr_d2 = []
for i in range(len(X)):
    arr_d2.append((dict_X_values_and_rangs[X[i]] - dict_Y_values_and_rangs[Y[i]])**2)
print(f"Массив квадратов разности рангов:")
for i in range(len(X)):
    print(f"d^2_{i+1} = {arr_d2[i]}")
print("Посчитаем коэф. корреляции Спирмана:")
r_s = 1 - 6*sum(arr_d2)/(len(X)**3 - len(X))
print(f"r_s = {r_s}")
print("Проверка значимости осуществляется по Стьюденту со степенью свободы df = n-2")
print(f"t_набл = {r_s * math.sqrt(len(X)-2) / math.sqrt(1-r_s**2)}")
print(f"Определение значимости и нахождение доверительного интервала оставляю на тебя, мой друг")
print("---")
print("Найдем коэф. корреляции Кендалла:")
print("Для этого изменим одинаковые ранги для выборок X, Y:")
rangs_X = rangs_for_Kendal(X.copy())
rangs_Y = rangs_for_Kendal(Y.copy())
scheme = []
for i in range(len(X)):
    scheme.append((rangs_X[i], rangs_Y[i]))
scheme.sort()
count_of_inversion = 0
for i in range(len(X)):
    for j in range(i+1, len(X)):
        if(scheme[i][1] > scheme[j][1]):
            count_of_inversion += 1
r_k = 1 - 4*count_of_inversion/(len(X)*(len(X)-1))
print(f"(rangs for X, rangs for Y): {scheme}")
print(f"Количество инверсий Q = {count_of_inversion}")
print(f"r_k = 1 - 4*Q/(n*(n-1))")
print(f"r_k = {r_k}")
print("Определим значимость:")
z = math.sqrt(9*len(X)*(len(X)-1))*r_k/math.sqrt(2*(2*len(X)+5))
if math.fabs(z) < 1.96:
    print(f"z = {z} и так как z < 1.96 -> гипотеза принимется о незначимости")
else:
    print(f"z = {z} и так как z >= 1.96 -> гипотеза непринимется о незначимости")
print()

print(f"Решим пункт 4:")
print("Линейное уравнение регрессии")
print(f"Решим систему: k{sum2(X)} + b{sum(X)} = {sumXY(X,Y)}\n\
               k{sum(X)} + b{len(X)} = {sum(Y)}\n")
systemArr = np.array(((sum2(X), sum(X)), (sum(X), len(X))))
VectorArr = np.array((sumXY(X,Y), sum(Y)))
solve = np.linalg.solve(systemArr, VectorArr)
print(f"k = {solve[0]}, b = {solve[1]}\n")
E = []
for i in range(len(X)):
    E.append(Y[i] - (solve[0]*X[i] + solve[1]))
    print(f"E_{i+1} = {E[i]}")
print()
print(f"1/nE e_i^2 = {sum2(E)/len(X)}")
print(f"R^2 = {1 - sum2(E)/(len(X)*D(Y))}")
R2 = 1 - sum2(E)/(len(X)*D(Y))
print("Найдем F_набл для статистики критерия о незначимости линейного уравнения регрессии:")
print(f"F_набл = {R2*(len(X) - 2)/(1 - R2)}")
print(f"F_т = F_0.05(1, {len(X) - 2})")

print("---\nКвадратное уравнение регрессии")
print(f"Решим систему: {sumN(X, 4)}a_2 + {sumN(X, 3)}a_1 + {sumN(X, 2)}a_0 = {sumXNY(X,Y,2)}\n\
               {sumN(X, 3)}a_2 + {sumN(X, 2)}a_1 + {sumN(X, 1)}a_0 = {sumXNY(X,Y,1)}\n\
               {sumN(X, 2)}a_2 + {sumN(X, 1)}a_1 + {len(X)}a_0 = {sum(Y)}\n")
systemArr = np.array(((sumN(X, 4), sumN(X, 3), sumN(X, 2)), (sumN(X, 3), sumN(X, 2), sumN(X, 1)), (sumN(X, 2), sumN(X, 1), len(X))))
VectorArr = np.array((sumXNY(X,Y,2), sumXNY(X,Y,1), sum(Y)))
solve = np.linalg.solve(systemArr, VectorArr)
print(f"a_2 = {solve[0]}, a_1 = {solve[1]}, a_0 = {solve[2]}\n")
E = []
for i in range(len(X)):
    E.append(Y[i] - (solve[0]*X[i]**2 + solve[1]*X[i] + solve[2]))
    print(f"E_{i+1} = {E[i]}")
print()
print(f"1/nE e_i^2 = {sum2(E)/len(X)}")
print(f"R^2 = {1 - sum2(E)/(len(X)*D(Y))}")
R2 = 1 - sum2(E)/(len(X)*D(Y))
print("Найдем F_набл для статистики критерия о незначимости квадратного уравнения регрессии:")
print(f"F_набл = {R2*(len(X) - 3)/((1 - R2)*2)}")
print(f"F_т = F_0.05(2, {len(X) - 3})")

print("---\nКубическое уравнение регрессии")
print(f"Решим систему: {sumN(X,6)}a_3 + {sumN(X, 5)}a_2 + {sumN(X, 4)}a_1 + {sumN(X, 3)}a_0 = {sumXNY(X,Y,3)}\n\
               {sumN(X,5)}a_3 + {sumN(X, 4)}a_2 + {sumN(X, 3)}a_1 + {sumN(X, 2)}a_0 = {sumXNY(X,Y,2)}\n\
               {sumN(X,4)}a_3 + {sumN(X, 3)}a_2 + {sumN(X, 2)}a_1 + {sumN(X, 1)}a_0 = {sumXNY(X,Y,1)}\n\
               {sumN(X,3)}a_3 + {sumN(X, 2)}a_2 + {sumN(X, 1)}a_1 + {len(X)}a_0 = {sum(Y)}\n")
systemArr = np.array(((sumN(X, 6), sumN(X, 5), sumN(X, 4), sumN(X, 3)), (sumN(X, 5), sumN(X, 4), sumN(X, 3), sumN(X, 2)), 
                      (sumN(X, 4), sumN(X, 3), sumN(X, 2), sumN(X, 1)), (sumN(X, 3), sumN(X, 2), sumN(X, 1), len(X) )))
VectorArr = np.array((sumXNY(X,Y,3), sumXNY(X,Y,2), sumXNY(X,Y,1), sum(Y)))
solve = np.linalg.solve(systemArr, VectorArr)
print(f"a_3 = {solve[0]}, a_2 = {solve[1]}, a_1 = {solve[2]}, a_0 = {solve[3]}\n")
E = []
for i in range(len(X)):
    E.append(Y[i] - (solve[0]*X[i]**3 + solve[1]*X[i]**2 + solve[2]*X[i] + solve[3]))
    print(f"E_{i+1} = {E[i]}")
print()
print(f"1/nE e_i^2 = {sum2(E)/len(X)}")
print(f"R^2 = {1 - sum2(E)/(len(X)*D(Y))}")
R2 = 1 - sum2(E)/(len(X)*D(Y))
print("Найдем F_набл для статистики критерия о незначимости кубического уравнения регрессии:")
print(f"F_набл = {R2*(len(X) - 3)/((1 - R2)*3)}")
print(f"F_т = F_0.05(2, {len(X) - 4})")

print("\nРешение пункта 5:")
print("Показательная регрессия")
ln_Y = []
for y in Y:
    ln_Y.append(math.log(y))

print(f"ln_Y = {ln_Y}")
print(f"M(ln_Y) = {M(ln_Y)}, D(ln_Y) = {D(ln_Y)}, cov(X, ln_Y) = {cov(X, ln_Y)}")
print(f"k = {k(X, ln_Y)}, b = {b(X, ln_Y)}")
print(f"a^ = {math.exp(k(X, ln_Y))}, b^ = {math.exp(b(X, ln_Y))}")

E = []
print()
for i in range(len(X)):
    E.append(Y[i] - (math.exp(b(X, ln_Y)) * (math.exp(k(X, ln_Y))**X[i])))
    print(f"E_{i+1} = {E[i]}")

print(f"1/nE e_i^2 = {sum2(E)/len(X)}")

print("\nСтепенная регрессия")
ln_X = []
for x in X:
    ln_X.append(math.log(x))

print(f"ln_X = {ln_X}")
print(f"M(ln_X) = {M(ln_X)}, D(ln_X) = {D(ln_X)}, cov(ln_X, ln_Y) = {cov(ln_X, ln_Y)}")
print(f"k = {k(ln_X, ln_Y)}, b = {b(ln_X, ln_Y)}")
print(f"a^ = {k(ln_X, ln_Y)}, b^ = {math.exp(b(ln_X, ln_Y))}")

E = []
print()
for i in range(len(X)):
    E.append(Y[i] - (math.exp(b(ln_X, ln_Y)) * (X[i] ** k(ln_X, ln_Y))))
    print(f"E_{i+1} = {E[i]}")

print(f"1/nE e_i^2 = {sum2(E)/len(X)}")

print("\nОбратно-линейная регресия")
Y_1 = []
for y in Y:
    Y_1.append(1/(y))

print(f"Y_1 = {Y_1}")
print(f"M(Y_1) = {M(Y_1)}, D(Y_1) = {D(Y_1)}, cov(X, Y_1) = {cov(X, Y_1)}")
print(f"k = {k(X, Y_1)}, b = {b(X, Y_1)}")
print(f"a^ = {k(X, Y_1)}, b^ = {b(X, Y_1)}")

E = []
print()
for i in range(len(X)):
    E.append(Y[i] - 1/(k(X, Y_1)*X[i] + b(X, Y_1)))
    print(f"E_{i+1} = {E[i]}")

print(f"1/nE e_i^2 = {sum2(E)/len(X)}")
input()
