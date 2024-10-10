# findlevel
指定したディレクトリ内のBMSファイルのレベル情報(PLAYLEVEL,DIFFICULTY)を収集し、CSVファイルに出力します。

- DIFFICULTYが指定されていない場合はログに出力されます。

## Usage
```
findlevel [dirpath]
```

## Example
```
> findlevel C:\BMS\BOFNT
Difficulty is missing: C:\BMS\BOFNT\[256mochi] TYFYHW\TYFYHW[Hyper].bms, #DIFFICULTY = 
Difficulty is missing: C:\BMS\BOFNT\[287 a.k.a. EXAM.S&在 Feat.アイノユメ] Scream_Our_Instinct!!\Scream_Our_Instinct!!(Another).bms, #DIFFICULTY =
...
Difficulty is missing, but this has only one chart: C:\BMS\BOFNT\[Cxx] Retnab\Retnab.bms, #DIFFICULTY =
Difficulty is missing, but this has only one chart: C:\BMS\BOFNT\[DJ Grief] Mystic Valley\1.bms, #DIFFICULTY =
...
Difficulties OK: 2230, NG: 178
Output file created: findlevel_output.csv
```

findlevel_output.csv
```
directory_path,song_name,for_aviutl,level_text
[0 K] VOYAGER III,VOYAGER III,"7KEY <p64,+0><#00ccff>☆3 <#f6962b>☆7 <#ff3014>☆9 <#>",7K: N3 H7 A9 
[0-n] Gingeeer Aaale !,Gingeeer Aaale !,"7KEY <p64,+0><#24fb5d>☆3 <#00ccff>☆6 <#f6962b>☆10 
<p64,+0><#ff3014>☆12 <#a049ff>☆9 <#>
14KEY <p76,+0><#24fb5d>☆3 <#00ccff>☆8 <#f6962b>☆10 
<p76,+0><#ff3014>☆12 <#>
9KEY <p64,+0><#00ccff>Lv24 <#f6962b>Lv43 <#ff3014>Lv48 
<p64,+0><#a049ff>Lv50 <#>","7K: B3 N6 H10 A12 I9 
14K: B3 N8 H10 A12 
9K: N24 H43 A48 I50 "
[14.14] InterConnected,InterConnected,"7KEY <p64,+0><#24fb5d>☆1 <#00ccff>☆3 <#f6962b>☆6 
<p64,+0><#ff3014>☆9 <#>",7K: B1 N3 H6 A9 
...
[黒皇帝 feat. AKA] Belle de Nuit,Belle de Nuit,"7KEY <p64,+0><#24fb5d>☆1 <#00ccff>☆4 <#f6962b>☆7 
<p64,+0><#ff3014>☆10 <#a049ff>☆12 <#>",7K: B1 N4 H7 A10 I12 
```