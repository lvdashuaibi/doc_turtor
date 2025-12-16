/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/callbacks"
	ecmodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// 颜色代码常量
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

func Event(event *adk.AgentEvent) {
	fmt.Printf("name: %s\npath: %s", event.AgentName, event.RunPath)
	if event.Output != nil && event.Output.MessageOutput != nil {
		if m := event.Output.MessageOutput.Message; m != nil {
			if len(m.Content) > 0 {
				if m.Role == schema.Tool {
					fmt.Printf("\ntool response: %s", m.Content)
				} else {
					fmt.Printf("\nanswer: %s", m.Content)
				}
			}
			if len(m.ToolCalls) > 0 {
				for _, tc := range m.ToolCalls {
					fmt.Printf("\ntool name: %s", tc.Function.Name)
					fmt.Printf("\narguments: %s", tc.Function.Arguments)
				}
			}
		} else if s := event.Output.MessageOutput.MessageStream; s != nil {
			toolMap := map[int][]*schema.Message{}
			var contentStart bool
			charNumOfOneRow := 0
			maxCharNumOfOneRow := 120
			for {
				chunk, err := s.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Printf("error: %v", err)
					return
				}
				if chunk.Content != "" {
					if !contentStart {
						contentStart = true
						if chunk.Role == schema.Tool {
							fmt.Printf("\ntool response: ")
						} else {
							fmt.Printf("\nanswer: ")
						}
					}

					charNumOfOneRow += len(chunk.Content)
					if strings.Contains(chunk.Content, "\n") {
						charNumOfOneRow = 0
					} else if charNumOfOneRow >= maxCharNumOfOneRow {
						fmt.Printf("\n")
						charNumOfOneRow = 0
					}
					fmt.Printf("%v", chunk.Content)
				}

				if len(chunk.ToolCalls) > 0 {
					for _, tc := range chunk.ToolCalls {
						index := tc.Index
						if index == nil {
							log.Fatalf("index is nil")
						}
						toolMap[*index] = append(toolMap[*index], &schema.Message{
							Role: chunk.Role,
							ToolCalls: []schema.ToolCall{
								{
									ID:    tc.ID,
									Type:  tc.Type,
									Index: tc.Index,
									Function: schema.FunctionCall{
										Name:      tc.Function.Name,
										Arguments: tc.Function.Arguments,
									},
								},
							},
						})
					}
				}
			}

			for _, msgs := range toolMap {
				m, err := schema.ConcatMessages(msgs)
				if err != nil {
					log.Fatalf("ConcatMessage failed: %v", err)
					return
				}
				fmt.Printf("\ntool name: %s", m.ToolCalls[0].Function.Name)
				fmt.Printf("\narguments: %s", m.ToolCalls[0].Function.Arguments)
			}
		}
	}
	if event.Action != nil {
		if event.Action.TransferToAgent != nil {
			fmt.Printf("\naction: transfer to %v", event.Action.TransferToAgent.DestAgentName)
		}
		if event.Action.Interrupted != nil {
			for _, ic := range event.Action.Interrupted.InterruptContexts {
				str, ok := ic.Info.(fmt.Stringer)
				if ok {
					fmt.Printf("\n%s", str.String())
				} else {
					fmt.Printf("\n%v", ic.Info)
				}
			}
		}
		if event.Action.Exit {
			fmt.Printf("\naction: exit")
		}
	}
	if event.Err != nil {
		fmt.Printf("\nerror: %v", event.Err)
	}
	fmt.Println()
	fmt.Println()
}

// NewOutputCallbackHandler 创建一个新的输出回调处理器，用于打印 LLM 和 Agent 的输出结果
func NewOutputCallbackHandler() callbacks.Handler {
	builder := callbacks.NewHandlerBuilder()

	// 处理开始事件
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		if info.Component == "ChatModel" {
			fmt.Printf("%s[%s] 开始调用 LLM: %s%s\n", Cyan, info.Type, info.Name, Reset)
		} else if info.Component == "Tool" {
			toolInput := tool.ConvCallbackInput(input)
			fmt.Printf("%s[工具调用] 开始执行工具: %s%s\n", Yellow, info.Name, Reset)
			fmt.Printf("%s参数: %s%s\n", Yellow, toolInput.ArgumentsInJSON, Reset)
		}
		return ctx
	})

	// 处理结束事件
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		if info.Component == "ChatModel" {
			modelOutput := ecmodel.ConvCallbackOutput(output)
			if modelOutput.Message != nil {
				fmt.Printf("%s[%s] LLM 输出 (%s):%s\n", Green, info.Type, info.Name, Reset)
				fmt.Printf("%s角色: %s%s\n", Green, modelOutput.Message.Role, Reset)
				fmt.Printf("%s内容: %s%s\n", Green, modelOutput.Message.Content, Reset)
				if len(modelOutput.Message.ToolCalls) > 0 {
					fmt.Printf("%s工具调用数: %d%s\n", Green, len(modelOutput.Message.ToolCalls), Reset)
					for i, tc := range modelOutput.Message.ToolCalls {
						fmt.Printf("%s  [%d] %s: %s%s\n", Green, i+1, tc.Function.Name, tc.Function.Arguments, Reset)
					}
				}
			}
		} else if info.Component == "Tool" {
			toolOutput := tool.ConvCallbackOutput(output)
			fmt.Printf("%s[工具结果] %s 执行完成:%s\n", Blue, info.Name, Reset)
			fmt.Printf("%s结果: %s%s\n", Blue, toolOutput.Response, Reset)
		}
		return ctx
	})

	// 处理流式输出
	builder.OnEndWithStreamOutputFn(func(ctx context.Context, info *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
		if info.Component == "ChatModel" {
			fmt.Printf("%s[%s] LLM 流式输出 (%s):%s\n", Magenta, info.Type, info.Name, Reset)

			go func() {
				defer output.Close()
				for {
					frame, err := output.Recv()
					if err != nil {
						if err == io.EOF {
							break
						}
						fmt.Printf("%s[错误] 接收流式输出失败: %v%s\n", Red, err, Reset)
						return
					}

					if modelOutput, ok := frame.(*ecmodel.CallbackOutput); ok {
						if modelOutput.Message != nil && modelOutput.Message.Content != "" {
							fmt.Printf("%s%s%s", Magenta, modelOutput.Message.Content, Reset)
						}
					}
				}
				fmt.Printf("\n%s[完成] 流式输出结束%s\n", Magenta, Reset)
			}()
		}
		return ctx
	})

	// 处理错误事件
	builder.OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
		fmt.Printf("%s[错误] %s (%s) 执行失败: %v%s\n", Red, info.Name, info.Component, err, Reset)
		return ctx
	})

	return builder.Build()
}

// NewDetailedOutputCallbackHandler 创建一个新的详细输出回调处理器（包括完整的 JSON）
func NewDetailedOutputCallbackHandler(writer io.Writer) callbacks.Handler {
	if writer == nil {
		writer = io.Discard
	}

	builder := callbacks.NewHandlerBuilder()

	// 处理开始事件
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Fprintf(writer, "%s[开始] %s (%s:%s)%s\n", Cyan, info.Name, info.Component, info.Type, Reset)
		if inputJSON, err := json.MarshalIndent(input, "  ", "  "); err == nil {
			fmt.Fprintf(writer, "%s输入: %s%s\n", Cyan, string(inputJSON), Reset)
		}
		return ctx
	})

	// 处理结束事件
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Fprintf(writer, "%s[完成] %s (%s:%s)%s\n", Green, info.Name, info.Component, info.Type, Reset)
		if outputJSON, err := json.MarshalIndent(output, "  ", "  "); err == nil {
			fmt.Fprintf(writer, "%s输出: %s%s\n", Green, string(outputJSON), Reset)
		}
		return ctx
	})

	// 处理错误事件
	builder.OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
		fmt.Fprintf(writer, "%s[错误] %s (%s:%s): %v%s\n", Red, info.Name, info.Component, info.Type, err, Reset)
		return ctx
	})

	return builder.Build()
}

// NewSimpleOutputCallbackHandler 创建一个新的简洁输出回调处理器
func NewSimpleOutputCallbackHandler() callbacks.Handler {
	builder := callbacks.NewHandlerBuilder()

	// 处理开始事件
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("%s[%s] 开始: %s%s\n", Cyan, info.Component, info.Name, Reset)
		return ctx
	})

	// 处理结束事件
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("%s[%s] 完成: %s%s\n", Green, info.Component, info.Name, Reset)
		return ctx
	})

	// 处理错误事件
	builder.OnErrorFn(func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
		fmt.Printf("%s[%s] 错误: %s - %v%s\n", Red, info.Component, info.Name, err, Reset)
		return ctx
	})

	return builder.Build()
}
