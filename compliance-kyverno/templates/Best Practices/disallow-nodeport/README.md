# Disallow Nodeports

Regla que impide que los svc usen NodePort.

## Descripción

Es una mala practica usar Nodeport, ya que expone un puerto a nivel de nodo para un pod específico, debe usarse Ingress, Routes o Svc.

## Parametros

No recibe parámetros

## Recursos

Por defecto actua sobre SVC unicamente

## Filtros

No recibe filtros
