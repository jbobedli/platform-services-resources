# Disallow Latest Tag

Regla que impide que los pods usen containers con tag :latest.

## Descripción

Utilizar la tag latest es una mala practica porque dicha etiqueta se refiere a una versión cambiante del código, lo que puede llevar a escenarios donde no se controla qué versión de código está presente dentro del Pod.

Esta norma obliga a utilizar una versión concreta del código.

## Parametros

No recibe parámetros

## Recursos

Por defecto actua sobre Pods unicamente

## Filtros

No recibe filtros
