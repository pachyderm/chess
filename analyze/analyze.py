import pandas as pd
import json
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import os

indir = '/pfs/annotate'

for root, dirs, filenames in os.walk(indir):
        
    for game_file in filenames:
        
        white = []
        black = []
        result = []
        
        with open('/pfs/annotate/' + game_file) as analyzed_game:
            
            white_move = 1
            black_move = 1
            overall_move_num = 1
        
            for move in analyzed_game:
                
                parsed_move = json.loads(move)
                
                score_diff = parsed_move["played-move-score"] - parsed_move["best-move-score"]

                if parsed_move["mover"] == parsed_move["white"]:
                    white.append([white_move, score_diff])
                    white_move += 1
                else:
                    black.append([black_move, score_diff])
                    black_move += 1

                threshold = 20

                if parsed_move["played-move-score"] == 0:
                    result.append([overall_move_num, 0.5])
                    overall_move_num += 1
                elif parsed_move["played-move-score"] > threshold and parsed_move["mover"] == parsed_move["white"]:
                    result.append([overall_move_num, 1])
                    overall_move_num += 1
                elif parsed_move["played-move-score"] < -threshold and parsed_move["mover"] == parsed_move["white"]:
                    result.append([overall_move_num, 0])
                    overall_move_num += 1
                elif parsed_move["played-move-score"] > threshold and parsed_move["mover"] == parsed_move["black"]:
                    result.append([overall_move_num, 0])
                    overall_move_num += 1
                elif parsed_move["played-move-score"] < -threshold and parsed_move["mover"] == parsed_move["black"]:
                    result.append([overall_move_num, 1])
                    overall_move_num += 1

        # Plot the results.
        whiteDF = pd.DataFrame(white, columns=["move", "score_diff"])
        blackDF = pd.DataFrame(black, columns=["move", "score_diff"])
        resultDF = pd.DataFrame(result, columns=["move", "result"])

        f, axarr = plt.subplots(3, figsize=(12,10))

        axarr[0].plot(whiteDF["move"], whiteDF["score_diff"])
        axarr[0].set_title('White - ' + parsed_move["white"])
        axarr[0].set_xlabel("White move number")
        axarr[0].set_ylabel("score delta (centipawns)")
        axarr[0].set_ylim([-200, 20])

        axarr[1].plot(blackDF["move"], blackDF["score_diff"])
        axarr[1].set_title('Black - ' + parsed_move["black"])
        axarr[1].set_xlabel("Black move number")
        axarr[1].set_ylabel("score delta (centipawns)")
        axarr[1].set_ylim([-200, 20])

        axarr[2].plot(resultDF["move"], resultDF["result"])
        axarr[2].set_title('Winning')
        axarr[2].set_xlabel("Overall move number")
        axarr[2].set_ylim([-0.2, 1.2])
        axarr[2].set_yticks([0, 0.5, 1])
        axarr[2].set_yticklabels(["Black", "Even", "White"])

        plt.tight_layout()
        plt.savefig('/pfs/out/' + game_file)

        # Output the score_diffs.
        whiteDF.to_csv("/pfs/out/white_" + game_file + ".csv")
        blackDF.to_csv("/pfs/out/black_" + game_file + ".csv")

