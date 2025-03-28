package estructuras

type SUBERBLOCK struct {
	s_filesystem_type   int32
	s_inodes_count      int32
	s_blocks_count      int32
	s_free_blocks_count int32
	s_free_inodes_count int32
	s_mtime             float32
	s_umtime            float32
	s_mnt_count         int32
	s_magic             int32
	s_inode_s           int32
	s_block_s           int32
	s_first_ino         int32
	s_first_blo         int32
	s_bm_inode_start    int32
	s_bm_block_start    int32
	s_inode_start       int32
	s_block_start       int32
}
