# CryptoMotionCoin la cryptomonaie meta
# aurelien592 la cryptomonaie meta

CryptoMotionCoin à pour objectif de permettre la création d'une cryptomonnaie innovante
aurelien592 a pour objectif de permettre la création d'une cryptomonnaie
à part entière.

Aujourd'hui, en 2018, il existe deux manières de créer une cryptomonnaie :
@@ -167,18 +167,135 @@ il pourrait même tenter de créer des scripts turing complet et produire des
contrats intelligents sur une base de cryptomonnaie similaire à Bitcoin, puisque
le protocole maître est proche de Bitcoin.

# Dictionnaire
# Algorithme de consensus

Intégration de l'algorithme AlgoRand.

AlgoRand est un algorithme de consenus proposé par Silvo Micali dans le cadre du
MIT. La proposition tente de palier aux problèmes de la preuve de travail de Bitcoin
qui sont la consommation en électricité élevée et le temps de validation trop long.

L'une des partcularités du protocole maître sera de recenser un dictionnaire des
protocoles annexes.
Algorand n'intègre pas de récompense pour le travail des acteurs, dans sa version actuelle.

Un dictionnaire se pésente sous la forme d'un script qui comporte toutes les informations
relatives à son protocole.
AlgoRand permet de produire un bloc toutes les 10 secondes sans que la sécurité soit
négligée. aurelien592 a une production de bloc à raison de 1 bloc par 1 minute.

Nous verrons que aurelien592 intègre un intérêt pour les acteurs à être actifs sur le réseau.
Cet aspect est crucial puisque la qualité du réseau réside dans le nombre d'acteurs
qui assurent la sécurité.

# Dictionnaire

De cette manière les noeuds sur le protocole maître seront notifiés de l'arrivée
d'un nouveau protocole et pourront le joindre si ils le souhaite. La topologie de réseau
de ce nouveau protocole annexe sera alors proche de la topologie du réseau
maître.
L'une des partcularités du protocole maître est de recenser un dictionnaire des
protocoles annexes. Ce recensemment se fait par le biais de transactions correspondant
au dictionnaire, exactement de la même manière qu'opère Ethereum pour les transactions
qui concernent un contrat (création ou appel d'une méthode). Une entrée de dictionnaire
est alors une adresse publique.

## Entrée de dictionnaire (protocole annexe publié)

Une entrée inclut les informations suivantes :
- Metadonnées
  - Nom de jeton
  - Nom du créateur
  - Adresse publique du créateur
  - Essence
  - Maximum d'acteurs souhaité
  - % de redistribution par tour de participation
  - Description _optionnel_
- Bloc code
- Transaction code
- Commandes code
- OPCodes code
- Consensus algorithme code

Une entrée est créée par une transaction.

Lorsqu'un créateur publie son protocole annexe il doit inclure un montant RZM qu'on
appelle le carburant. Ce montant sert de récompense aux acteurs du protocole
annexe, expliqué en détail plus bas. Concrètement les jetons inclus ont pour destin
d'être détruits, la redistribution se fait sous forme de _coinbase_ comme le système
de récompense dans Bitcoin.

## Intérêt pour les acteurs

Le protocole annexe, selon sa conception produit également des jetons qui lui sont propres,
cela peut justifier d'un intérêt à être acteur sur ce protocole en dehors de aurelien592.
En d'autres termes, un protocole annexe peut perdurer sans dépendre du réseau
aurelien592, avec sa propre économie et ses acteurs.

Les acteurs sur le protocole annexe doivent prouver leur activité sur ce dernier, afin
d'être élligibles à la consommation des fonds de l'entrée du dictionnaire, recevoir
la coinbase. Etant donné que l'algorithme de consensus est défini par le créateur du protocole
annexe, il n'est alors pas possible de générer une élligibilité sans confiance (trustless)
à partir du protocole annexe. Les parties suivantes détaillent comment aurelien592
pallie à cette problématique.

### Tour de participation

Un tour de participation est une période où tous les acteurs déclarent participer
à x protocole. ils renseignent leur adresse IP avec leur adresse publique aurelien592,
ces informations sont recensées lors de la validation du tour de participation.
Un tour de participation par entrée de dictionnaire. Un tour de participation débute
tous les 59 blocs, après le tour précédent, ce qui correspond à 59 minutes.
Cela signifie qu'un protocole annexe, lorsqu'il est publié, n'est considéré
par le réseau qu'après 59 blocs qui suit sa publication. Cela laisse le temps
aux acteurs de le notifier et de s'inscrire comme participants au prochain tour.

### Déterminer l'activité des acteurs en fin de tour de participation

Avant le début d'un nouveau tour de participation un quorum d'acteurs est choisi
aléatoirement. Son rôle est de contrôler l'activité des participants.
L'activité est contrôlée en réalisant un `PING` vers les participants.
Un `PING` est lancé toutes les minutes et la réponse est recensée par chaque
acteur du quorum indépendemment. Une réponse compte 1 point pour le participant,
une absence de réponse apporte aucun point. En fin de tour de participation,
tous les acteurs du quorum envoient leur rapport aux validateurs.

Les validateurs pondèrent les points de chaque participant du tour afin de déterminer
quels participants ont été suffisamment actifs durant ce tour de participation.
On détermine l'activité suffisante selon le pseudocode suivant :

```
algorithme determiner-activite est
    entrée: Des rapports R
    sortie: Liste des participants déterminés actifs L
    On suppose P une liste vide telle que <adresse, points> où l'adresse
    est l'adresse publique du participant
    tant qu'il existe un rapport r dans R faire
        pour chaque entrée (adresse, points) dans r faire
            P[adresse] += points
    On suppose L une liste vide
    On suppose B := 59 qui correspond au nombre de point total maximal sur un rapport
    On support nb_R le nombre de rapport
    pour chaque entrée (adresse, points) dans P faire
        si (points / nb_R) > (B x 0.8)
            Ajouter adresse à L
    retourner L
```

On considère qu'un acteur a été suffisamment actif à +80% de réponse durant le
tour de participation.

Tous les validateurs doivent arriver à la même liste de participants, ainsi le consensus
peut-être atteint et la _coinbase_ peut être ditribuée à tous les participants
déterminés comme actifs durant le tour de participation. Le montant de la _coinbase_
est définie ainsi :

coinbase = (essence * (redistribution / 100)) / nb_participants

Où
- essence est la quantité de jeton injecté dans le dictionnaire au début du tour de
participation
- redistribution est le % de redistribution indiqué par le dictionnaire au début
du tour de participation
- nb_participants est le nombre de participant déterminés comme actifs

# Conclusion

Ce projet propose un échaffaudage à cryptomonnaies. De nombreuses personnes ressentent le besoin
de faire leur cryptomonnaie pour des besoins particuliers ou même pour s'amuser à
avoir sa cryptomonnaie.
En plus de l'utilité économique au même titre que Bitcoin, avec l'ancienneté en moins,
ce projet s'adresse également à l'académie et à la recherche qui pourrait avoir
besoin de mettre à l'épreuve des algorithmes de consensus. Cette cryptomonnaie
est un gain de temps considérable dans la mesure où elle propose de construire
sa cryptomonnaie et de bénéficier d'une topologie de réseau déjà acquise.