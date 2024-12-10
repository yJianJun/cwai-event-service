package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SortGeneratorTestSuite struct {
	suite.Suite
	generator *SortGenerator
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(SortGeneratorTestSuite))
}

type TestEmbeddingStruct struct {
	CPU int32 `json:"cpu" gorm:"column:cpu_total;type:int"`
	GPU int32 `json:"gpu" gorm:"type:int;column:gpu_total"`
	MEM int32 `json:"mem" gorm:"column:mem_total;type:int"`
}

type TestStruct struct {
	TestEmbeddingStruct
	CPUUtil   float64 `json:"cpuUtil"`
	GPUUtil   float64 `json:"gpuUtil"`
	MEMUtil   float64 `json:"memUtil"`
	StartTime int64   `gorm:"column:start_time" json:"startTime"`
}

func (suite *SortGeneratorTestSuite) SetupSuite() {
	suite.generator = NewSortGenerator(TestStruct{})
}

func (suite *SortGeneratorTestSuite) BeforeTesting() {
	assert.NotNil(suite.T(), suite.generator)
}

func (suite *SortGeneratorTestSuite) TestMapping() {
	t := suite.T()
	m := suite.generator.Mapping
	assert.Equal(t, "cpu_total", m["cpu"])
	assert.Equal(t, "mem_util", m["memUtil"])
	assert.Equal(t, "gpu_total", m["gpu"])
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_Empty() {
	t := suite.T()
	sql, e := suite.generator.Generate("")
	assert.NoError(t, e)
	assert.Empty(t, sql)
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_DBTag() {
	t := suite.T()
	sql, e := suite.generator.Generate("start_time,desc")
	assert.NoError(t, e)
	assert.Equal(t, "ORDER BY start_time DESC", sql)
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_Simple() {
	t := suite.T()
	sql, e := suite.generator.Generate("cpuUtil")
	assert.NoError(t, e)
	assert.Equal(t, "ORDER BY cpu_util", sql)
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_Simple_Descending() {
	t := suite.T()
	sql, e := suite.generator.Generate("gpu,desc")
	assert.NoError(t, e)
	assert.Equal(t, "ORDER BY gpu_total DESC", sql)
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_Multiple_Simple() {
	t := suite.T()
	sql, e := suite.generator.Generate("cpu;gpu;mem")
	assert.NoError(t, e)
	assert.Equal(t, "ORDER BY cpu_total, gpu_total, mem_total", sql)
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_Multiple_Descending() {
	t := suite.T()
	sql, e := suite.generator.Generate("cpu,desc;gpu,desc;mem,desc")
	assert.NoError(t, e)
	assert.Equal(t, "ORDER BY cpu_total DESC, gpu_total DESC, mem_total DESC", sql)
}

func (suite *SortGeneratorTestSuite) TestGenerateSQL_Multiple_Mixed() {
	t := suite.T()
	sql, e := suite.generator.Generate("cpu,asc;gpu,desc;mem")
	assert.NoError(t, e)
	assert.Equal(t, "ORDER BY cpu_total ASC, gpu_total DESC, mem_total", sql)
}

func (suite *SortGeneratorTestSuite) TestPtr() {
	g := NewSortGenerator(&TestStruct{})
	suite.Require().NotNil(g)
}

func (suite *SortGeneratorTestSuite) TestFactory() {
	g1, err := factory.GetGenerator(TestEmbeddingStruct{})
	suite.Nil(err)
	suite.NotNil(g1)

	g2, err := factory.GetGenerator(TestStruct{})
	suite.Nil(err)
	suite.NotNil(g2)

	g3, err := factory.GetGenerator(TestStruct{})
	suite.Nil(err)
	suite.NotNil(g3)
	suite.Equal(g2, g3)
	suite.NotEqual(g1, g3)
}

func (suite *SortGeneratorTestSuite) TestFactoryWithPointer() {
	g1, err := factory.GetGenerator(&TestEmbeddingStruct{})
	suite.Nil(err)
	suite.NotNil(g1)

	g2, err := factory.GetGenerator(&TestStruct{})
	suite.Nil(err)
	suite.NotNil(g2)

	g3, err := factory.GetGenerator(&TestStruct{})
	suite.Nil(err)
	suite.NotNil(g3)
	suite.Equal(g2, g3)
	suite.NotEqual(g1, g3)
}
