package main

import "testing"

//TestAdd は Add関数の動作を確認するテストです
func TestAdd(t *testing.T) {
	//1.期待する結果　(Expected)
	expected := 3

	//2.実際の実行結果　(Actual)
	actual := Add(1, 2)

	//3.答え合わせ
	if actual != expected {
		t.Errorf("Add( 1,2 ) failed. Expected %d, got %d", expected, actual)
	}
}
