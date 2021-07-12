set -eu
./eval --problem-id 1 --pose-file test/1_sample | grep "valid"
echo "OK"
./eval --problem-id 2 --pose-file test/2_invalid.solution | grep "invalid edge"
echo "OK"
./eval --problem-id 3 --pose-file test/3_invalid.solution | grep "invalid edge"
echo "OK"
./eval --problem-id 3 --pose-file test/3_out_vertex | grep  "out of the hole"
echo "OK"
./eval --problem-id 3 --pose-file test/3_outof_edge | grep  "invalid edge"
echo "OK"
./eval --problem-id 3 --pose-file test/3_wall_hack | grep  "valid"
echo "OK"
./eval --problem-id 3 --pose-file test/3_super_flex | grep  "valid"
echo "OK"
./eval --problem-id 5 --pose-file test/5_global.solution | grep "valid"
echo "OK"
./eval --problem-id 5 --pose-file test/5_invalid_epsilon | grep "global epsilon budget exceeded."
echo "OK"
./eval --problem-id 14 --pose-file test/14_leg | grep "valid"
echo "OK"
./eval --problem-id 59 --pose-file test/59_on_hole_edge | grep "valid"
echo "OK"
./eval --problem-id 68 --pose-file test/68_valid.solution | grep "valid"
echo "OK"
./eval --problem-id 105 --pose-file test/105_superflex | grep "valid"
echo "OK"
echo "END"