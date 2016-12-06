import pandas as pd
import json
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import seaborn as sns
import os

indir = '/pfs/analyze/'
gamedir = '/pfs/annotate/'

player1 = 'Sergey Karjakin'
player2 = 'Magnus Carlsen'

for root, dirs, filenames in os.walk(indir):
        
    results1 = {}
    results2 = {}

    results1_small = {}
    results2_small = {}

    for analysis_file in filenames:

        if analysis_file[-4:] == ".csv":
            
            color = analysis_file[0:5]
            game = int(analysis_file.split('.')[0][10:])

            dataDF = pd.read_csv(indir + analysis_file)

            blunders = len(dataDF[dataDF["score_diff"] <= -100])
            small = dataDF[dataDF["score_diff"] > -100]
            blunders_small = len(small[small["score_diff"] <= -50]) 

            with open(gamedir + 'game' + str(game)) as analyzed_game: 
                parsed_move = json.loads(analyzed_game.readline())

            player = parsed_move[color]

            if player == player1:
                results1[game] = blunders
                results1_small[game] = blunders_small
            elif player == player2:
                results2[game] = blunders
                results2_small[game] = blunders_small

last_game = max(results1.keys())
games = range(1, last_game+1)

resultsDF = pd.DataFrame(games, columns=['game'])
resultsDF[player1] = resultsDF['game'].map(lambda x: results1[x])
resultsDF[player2] = resultsDF['game'].map(lambda x: results2[x])

resultsSmallDF = pd.DataFrame(games, columns=['game'])
resultsSmallDF[player1] = resultsSmallDF['game'].map(lambda x: results1_small[x])
resultsSmallDF[player2] = resultsSmallDF['game'].map(lambda x: results2_small[x])

resultsDF.set_index('game', inplace=True)
resultsSmallDF.set_index('game', inplace=True)

ax = resultsDF.plot(kind='bar', rot=0)
ax.set_ylabel("1+ pawn blunders")
ax.set_xlabel("game (* indicates a rapid game)")
ax.set_xticks(range(0, last_game))
ax.set_xticklabels(['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', 
    '12', '*13', '*14', '*15', '*16'])
fig = ax.get_figure()
fig.savefig('/pfs/out/blunders.png')

ax = resultsSmallDF.plot(kind='bar', rot=0)
ax.set_ylabel("0.5-1 pawn blunders")
ax.set_xlabel("game (* indicates a rapid game)")
ax.set_xticks(range(0, last_game))
ax.set_xticklabels(['1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', 
    '12', '*13', '*14', '*15', '*16'])
fig = ax.get_figure()
fig.savefig('/pfs/out/blunders_small.png')
 

